package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformHealth "github.com/jefrryss/go-grpc-microservices/platform/pkg/grpc/health"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	inventoryV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	logger     *logger.Logger
	closer     *closer.Closer
	listener   net.Listener
	grpcServer *grpc.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg.Mongo.URI() == "" {
		return nil, errors.New("MONGO_URI is required")
	}

	log, err := logger.New(logger.Config{Level: cfg.Logger.Level(), JSON: cfg.Logger.JSON()})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(connectCtx, options.Client().ApplyURI(cfg.Mongo.URI()))
	if err != nil {
		return nil, fmt.Errorf("connect MongoDB: %w", err)
	}
	if err := mongoClient.Ping(connectCtx, readpref.Primary()); err != nil {
		_ = mongoClient.Disconnect(connectCtx)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Address())
	if err != nil {
		_ = mongoClient.Disconnect(connectCtx)
		return nil, fmt.Errorf("listen on %s: %w", cfg.GRPC.Address(), err)
	}

	container := newContainer(mongoClient.Database(cfg.Mongo.Database()), cfg.Mongo.Collection())
	grpcServer := grpc.NewServer()
	inventoryV1.RegisterInventoryServiceServer(grpcServer, container.API())
	platformHealth.Register(grpcServer, inventoryV1.InventoryService_ServiceDesc.ServiceName)
	reflection.Register(grpcServer)

	resourceCloser := closer.New()
	resourceCloser.Add("MongoDB", mongoClient.Disconnect)
	resourceCloser.Add("gRPC listener", func(context.Context) error {
		err := listener.Close()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	})
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })

	return &App{logger: log, closer: resourceCloser, listener: listener, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "InventoryService started")
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

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
