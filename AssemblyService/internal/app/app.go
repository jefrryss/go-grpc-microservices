package app

import (
	"context"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/messaging"
	"github.com/jefrryss/go-grpc-microservices/AssemblyService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformConsumer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/consumer"
	platformProducer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/producer"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	kafkaMiddleware "github.com/jefrryss/go-grpc-microservices/platform/pkg/middleware/kafka"
)

type App struct {
	logger   *logger.Logger
	consumer *messaging.Consumer
	closer   *closer.Closer
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

	publisher := messaging.NewPublisher(producer)
	assembly := service.NewAssembly(cfg.BuildDuration, publisher)
	return &App{logger: log, consumer: messaging.NewConsumer(consumer, assembly), closer: resourceCloser}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "AssemblyService started")
	if err := a.consumer.Run(ctx); err != nil {
		return fmt.Errorf("run AssemblyService: %w", err)
	}
	return nil
}

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
