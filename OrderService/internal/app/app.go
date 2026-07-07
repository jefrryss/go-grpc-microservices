package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/messaging"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformHealth "github.com/jefrryss/go-grpc-microservices/platform/pkg/grpc/health"
	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	platformConsumer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/consumer"
	platformProducer "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka/producer"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	kafkaMiddleware "github.com/jefrryss/go-grpc-microservices/platform/pkg/middleware/kafka"
	orderV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type App struct {
	logger     *logger.Logger
	closer     *closer.Closer
	listener   net.Listener
	grpcServer *grpc.Server
	httpServer *http.Server
	consumer   *messaging.ShipAssembledConsumer
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg.Database.URL() == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	log, err := logger.New(logger.Config{Level: cfg.Logger.Level(), JSON: cfg.Logger.JSON()})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	database, err := pgxpool.New(ctx, cfg.Database.URL())
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}
	if err := database.Ping(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	inventoryConn, err := grpc.NewClient(cfg.Dependencies.InventoryAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("create InventoryService client: %w", err)
	}
	paymentConn, err := grpc.NewClient(cfg.Dependencies.PaymentAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = inventoryConn.Close()
		database.Close()
		return nil, fmt.Errorf("create PaymentService client: %w", err)
	}

	listener, err := net.Listen("tcp", cfg.Server.GRPCAddress())
	if err != nil {
		_ = paymentConn.Close()
		_ = inventoryConn.Close()
		database.Close()
		return nil, fmt.Errorf("listen on %s: %w", cfg.Server.GRPCAddress(), err)
	}

	resourceCloser := closer.New()
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })
	resourceCloser.Add("PostgreSQL", func(context.Context) error { database.Close(); return nil })
	resourceCloser.Add("InventoryService connection", func(context.Context) error { return inventoryConn.Close() })
	resourceCloser.Add("PaymentService connection", func(context.Context) error { return paymentConn.Close() })
	resourceCloser.Add("gRPC listener", func(context.Context) error {
		err := listener.Close()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	})

	var orderPaidPublisher *messaging.OrderPaidPublisher
	var kafkaConsumer platformKafka.Consumer
	if len(cfg.Kafka.Brokers()) > 0 {
		producer, err := platformProducer.New(cfg.Kafka.Brokers(), cfg.Kafka.OrderPaidTopic(), log)
		if err != nil {
			_ = resourceCloser.Close(ctx)
			return nil, fmt.Errorf("create OrderPaid producer: %w", err)
		}
		resourceCloser.Add("OrderPaid producer", func(context.Context) error { return producer.Close() })
		orderPaidPublisher = messaging.NewOrderPaidPublisher(producer)

		consumer, err := platformConsumer.New(
			cfg.Kafka.Brokers(),
			cfg.Kafka.ShipAssembledGroupID(),
			[]string{cfg.Kafka.ShipAssembledTopic()},
			log,
			kafkaMiddleware.Logging(log),
		)
		if err != nil {
			_ = resourceCloser.Close(ctx)
			return nil, fmt.Errorf("create ShipAssembled consumer: %w", err)
		}
		resourceCloser.Add("ShipAssembled consumer", func(context.Context) error { return consumer.Close() })
		kafkaConsumer = consumer
	}

	container := newContainer(database, inventoryConn, paymentConn, orderPaidPublisher)
	var shipAssembledConsumer *messaging.ShipAssembledConsumer
	if kafkaConsumer != nil {
		shipAssembledConsumer = messaging.NewShipAssembledConsumer(kafkaConsumer, container.Service())
	}
	grpcServer := grpc.NewServer()
	orderV1.RegisterOrderServiceServer(grpcServer, container.API())
	platformHealth.Register(grpcServer, orderV1.OrderService_ServiceDesc.ServiceName)
	reflection.Register(grpcServer)

	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{UseProtoNames: true},
		}),
		runtime.WithForwardResponseOption(ForwardHTTPStatus),
	)
	if err := orderV1.RegisterOrderServiceHandlerFromEndpoint(
		ctx,
		gatewayMux,
		cfg.Server.GatewayEndpoint(),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return nil, fmt.Errorf("register gRPC gateway: %w", err)
	}

	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/redoc.html")
	})
	rootMux.HandleFunc("/swagger/order.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger/order.swagger.json")
	})
	rootMux.HandleFunc("/openapi/order.openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi/order.openapi.yaml")
	})
	rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
			return
		}
		gatewayMux.ServeHTTP(w, r)
	})
	httpServer := &http.Server{Addr: cfg.Server.HTTPAddress(), Handler: rootMux}

	resourceCloser.Add("gRPC server", func(context.Context) error { grpcServer.GracefulStop(); return nil })
	resourceCloser.Add("HTTP server", httpServer.Shutdown)

	return &App{
		logger: log, closer: resourceCloser, listener: listener,
		grpcServer: grpcServer, httpServer: httpServer, consumer: shipAssembledConsumer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "OrderService started")
	errCh := make(chan error, 3)
	go func() { errCh <- a.grpcServer.Serve(a.listener) }()
	go func() { errCh <- a.httpServer.ListenAndServe() }()
	if a.consumer != nil {
		go func() { errCh <- a.consumer.Run(ctx) }()
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err == nil {
			return nil
		}
		if errors.Is(err, grpc.ErrServerStopped) || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve OrderService: %w", err)
	}
}

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
