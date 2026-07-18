package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	platformMigrator "github.com/jefrryss/go-grpc-microservices/platform/pkg/migrator/pg"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New(logger.Config{Level: "info", JSON: true})
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := run(ctx); err != nil {
		log.Error(ctx, "Migration failed", zap.Error(err))
		_ = log.Sync()
		os.Exit(1)
	}
	log.Info(ctx, "Migrations applied successfully")
	_ = log.Sync()
}

func run(ctx context.Context) error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to PostgreSQL: %w", err)
	}

	if err := platformMigrator.New(db, migrationsPath).Up(ctx); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
