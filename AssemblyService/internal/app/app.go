package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/messaging"
	assemblyMetrics "github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/metrics"
	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformConsumer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/consumer"
	platformProducer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/producer"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	kafkaMiddleware "github.com/jefrryss/go-grpc-microservices/platform/pkg/middleware/kafka"
	platformPrometheus "github.com/jefrryss/go-grpc-microservices/platform/pkg/prometheus"
)

type App struct {
	logger   *logger.Logger
	consumer *messaging.Consumer
	closer   *closer.Closer
	metrics  *http.Server
}

func New(cfg *config.Config) (*App, error) {
	log, err := logger.New(logger.Config{Level: cfg.LoggerLevel, JSON: cfg.LoggerJSON})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}
	resourceCloser := closer.New()
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })

	producer, err := platformProducer.New(cfg.Brokers, cfg.ShipAssembledTopic, log)
	if err != nil {
		_ = resourceCloser.Close(context.Background())
		return nil, fmt.Errorf("create producer: %w", err)
	}
	resourceCloser.Add("producer", func(context.Context) error { return producer.Close() })

	consumer, err := platformConsumer.New(
		cfg.Brokers,
		cfg.GroupID,
		[]string{cfg.OrderPaidTopic},
		log,
		kafkaMiddleware.Logging(log),
	)
	if err != nil {
		_ = resourceCloser.Close(context.Background())
		return nil, fmt.Errorf("create consumer: %w", err)
	}
	resourceCloser.Add("consumer", func(context.Context) error { return consumer.Close() })

	registry := platformPrometheus.NewRegistry()
	metricsServer := &http.Server{Addr: cfg.MetricsAddress, Handler: registry.Handler()}
	resourceCloser.Add("metrics server", metricsServer.Shutdown)
	publisher := messaging.NewPublisher(producer)
	assembly := service.NewAssembly(cfg.BuildDuration, publisher, assemblyMetrics.NewAssembly(registry.Registerer()))
	return &App{
		logger: log, consumer: messaging.NewConsumer(consumer, assembly),
		closer: resourceCloser, metrics: metricsServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "AssemblyService started")
	errCh := make(chan error, 2)
	go func() { errCh <- a.consumer.Run(ctx) }()
	go func() { errCh <- a.metrics.ListenAndServe() }()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("run AssemblyService: %w", err)
	}
}

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
