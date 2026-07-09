package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/api"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/config"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/middleware"
	userRepository "github.com/jefrryss/go-grpc-microservices/AuthService/internal/repository/postgres"
	sessionRepository "github.com/jefrryss/go-grpc-microservices/AuthService/internal/repository/session"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/closer"
	platformHealth "github.com/jefrryss/go-grpc-microservices/platform/pkg/grpc/health"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	pgMigrator "github.com/jefrryss/go-grpc-microservices/platform/pkg/migrator/pg"
	authV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/auth/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	logger     *logger.Logger
	closer     *closer.Closer
	listener   net.Listener
	grpcServer *grpc.Server
	httpServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	log, err := logger.New(logger.Config{Level: cfg.LoggerLevel, JSON: cfg.LoggerJSON})
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	migrationDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open migration database: %w", err)
	}
	if err := pgMigrator.New(migrationDB, cfg.MigrationsPath).Up(ctx); err != nil {
		_ = migrationDB.Close()
		return nil, err
	}
	if err := migrationDB.Close(); err != nil {
		return nil, fmt.Errorf("close migration database: %w", err)
	}

	database, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}
	if err := database.Ping(ctx); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisAddress, Password: cfg.RedisPassword})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	sessions := sessionRepository.New(redisClient, cfg.SessionTTL)
	authService := service.New(userRepository.NewUserRepository(database), sessions)
	listener, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		_ = redisClient.Close()
		database.Close()
		return nil, fmt.Errorf("listen gRPC: %w", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(middleware.UnaryAuth(authService)))
	authV1.RegisterAuthServiceServer(grpcServer, api.New(authService))
	platformHealth.Register(grpcServer, authV1.AuthService_ServiceDesc.ServiceName)
	reflection.Register(grpcServer)

	gateway := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
		if strings.EqualFold(key, "session-uuid") {
			return "session-uuid", true
		}
		return runtime.DefaultHeaderMatcher(key)
	}))
	if err := authV1.RegisterAuthServiceHandlerFromEndpoint(
		ctx, gateway, cfg.GatewayTarget,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return nil, fmt.Errorf("register AuthService gateway: %w", err)
	}

	root := http.NewServeMux()
	root.HandleFunc("/authorize", authorizeHandler(authService))
	root.Handle("/", gateway)
	httpServer := &http.Server{Addr: cfg.HTTPAddress, Handler: root}

	resourceCloser := closer.New()
	resourceCloser.Add("logger", func(context.Context) error { return log.Sync() })
	resourceCloser.Add("PostgreSQL", func(context.Context) error { database.Close(); return nil })
	resourceCloser.Add("Redis", func(context.Context) error { return redisClient.Close() })
	resourceCloser.Add("gRPC listener", func(context.Context) error {
		err := listener.Close()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	})
	resourceCloser.Add("gRPC server", func(context.Context) error { grpcServer.GracefulStop(); return nil })
	resourceCloser.Add("HTTP server", httpServer.Shutdown)
	return &App{logger: log, closer: resourceCloser, listener: listener, grpcServer: grpcServer, httpServer: httpServer}, nil
}

func authorizeHandler(authService *service.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		sessionValue := request.Header.Get("session-uuid")
		if sessionValue == "" {
			sessionValue = strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
		}
		sessionID, err := uuid.Parse(sessionValue)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := authService.ResolveSession(request.Context(), sessionID)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Header().Set("x-user-uuid", userID.String())
		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info(ctx, "AuthService started")
	errCh := make(chan error, 2)
	go func() { errCh <- a.grpcServer.Serve(a.listener) }()
	go func() { errCh <- a.httpServer.ListenAndServe() }()
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if errors.Is(err, grpc.ErrServerStopped) || errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

func (a *App) Close(ctx context.Context) error { return a.closer.Close(ctx) }
