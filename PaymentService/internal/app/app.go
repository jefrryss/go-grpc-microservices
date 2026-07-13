package app

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformHealth "github.com/jefrryss/go-grpc-microservices/platform/pkg/grpc/health"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/tracing"
	paymentV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config     *config.Config
	logger     *logger.Logger
	closer     *closer.Closer
	listener   net.Listener
	grpcServer *grpc.Server
}

func New(cfg *config.Config) (*App, error) {
	log, err := logger.New(logger.Config{Level: cfg.Logger.Level(), JSON: cfg.Logger.JSON()})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}
	shutdownTracer, err := tracing.Init(context.Background(), tracing.Config{
		ServiceName: "payment-service", Endpoint: cfg.Tracing.Endpoint(),
	})
	if err != nil {
		return nil, fmt.Errorf("initialize tracing: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		return nil, fmt.Errorf("listen on %s: %w", cfg.GRPC.Address(), err)
	}

	container := newContainer(log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(tracing.UnaryServerInterceptor("payment-service")))
	paymentV1.RegisterPaymentServiceServer(grpcServer, container.API())
	platformHealth.Register(grpcServer, paymentV1.PaymentService_ServiceDesc.ServiceName)
	reflection.Register(grpcServer)

	resourceCloser := closer.New()
	resourceCloser.Add("gRPC listener", func(context.Context) error {
		err := listener.Close()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	})
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })
	resourceCloser.Add("tracing", shutdownTracer)

	return &App{config: cfg, logger: log, closer: resourceCloser, listener: listener, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "PaymentService started")
	errCh := make(chan error, 1)
	go func() { errCh <- a.grpcServer.Serve(a.listener) }()

	select {
	case <-ctx.Done():
		a.grpcServer.GracefulStop()
		return nil
	case err := <-errCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return fmt.Errorf("serve gRPC: %w", err)
	}
}

func (a *App) Close(ctx context.Context) error {
	return a.closer.Close(ctx)
}
