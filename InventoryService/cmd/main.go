package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/app"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/config"
)

func main() {
	if err := config.Load(os.Getenv("CONFIG_PATH")); err != nil {
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, config.AppConfig())
	if err != nil {
		panic(err)
	}

	if err := application.Run(ctx); err != nil {
		panic(err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.Close(shutdownCtx); err != nil {
		panic(err)
	}
}
