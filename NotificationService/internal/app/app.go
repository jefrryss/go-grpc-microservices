package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/client/auth"
	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/client/telegram"
	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/config"
	serviceConsumer "github.com/jefrryss/go-grpc-microservices/NotificationService/internal/consumer"
	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformConsumer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/consumer"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	kafkaMiddleware "github.com/jefrryss/go-grpc-microservices/platform/pkg/middleware/kafka"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	logger        *logger.Logger
	closer        *closer.Closer
	orderPaid     *serviceConsumer.Consumer
	shipAssembled *serviceConsumer.Consumer
	telegram      *telegram.Client
}

func New(cfg *config.Config) (*App, error) {
	log, err := logger.New(logger.Config{Level: cfg.LoggerLevel, JSON: cfg.LoggerJSON})
	if err != nil {
		return nil, err
	}
	resourceCloser := closer.New()
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })

	authConnection, err := grpc.NewClient(cfg.AuthAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("connect AuthService: %w", err)
	}
	resourceCloser.Add("AuthService connection", func(context.Context) error { return authConnection.Close() })

	orderConsumer, err := platformConsumer.New(
		cfg.Brokers, "notification-order-paid", []string{cfg.OrderPaidTopic}, log, kafkaMiddleware.Logging(log),
	)
	if err != nil {
		_ = resourceCloser.Close(context.Background())
		return nil, fmt.Errorf("create OrderPaid consumer: %w", err)
	}
	resourceCloser.Add("OrderPaid consumer", func(context.Context) error { return orderConsumer.Close() })

	assembledConsumer, err := platformConsumer.New(
		cfg.Brokers, "notification-ship-assembled", []string{cfg.ShipAssembledTopic}, log, kafkaMiddleware.Logging(log),
	)
	if err != nil {
		_ = resourceCloser.Close(context.Background())
		return nil, fmt.Errorf("create ShipAssembled consumer: %w", err)
	}
	resourceCloser.Add("ShipAssembled consumer", func(context.Context) error { return assembledConsumer.Close() })

	telegramClient := telegram.New(cfg.TelegramToken)
	notifier := service.NewNotify(auth.New(authConnection), telegramClient)
	return &App{
		logger: log, closer: resourceCloser, telegram: telegramClient,
		orderPaid:     serviceConsumer.NewOrderPaid(orderConsumer, notifier),
		shipAssembled: serviceConsumer.NewShipAssembled(assembledConsumer, notifier),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "NotificationService started")
	errCh := make(chan error, 3)
	go func() { errCh <- a.orderPaid.Run(ctx) }()
	go func() { errCh <- a.shipAssembled.Run(ctx) }()
	go func() { errCh <- a.telegram.Poll(ctx) }()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err == nil || errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
}

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
