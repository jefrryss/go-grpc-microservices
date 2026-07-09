package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/app"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/config"
)

func main() {
	cfg, err := config.Load(os.Getenv("CONFIG_PATH"))
	if err != nil {
		panic(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	application, err := app.New(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := application.Run(ctx); err != nil {
		panic(err)
	}
	closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.Close(closeCtx); err != nil {
		panic(err)
	}
}
