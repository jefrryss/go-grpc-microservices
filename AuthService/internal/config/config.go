package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jefrryss/go-grpc-microservices/platform/pkg/config/envfile"
)

type Config struct {
	DatabaseURL    string
	RedisAddress   string
	RedisPassword  string
	SessionTTL     time.Duration
	MigrationsPath string
	GRPCAddress    string
	GatewayTarget  string
	HTTPAddress    string
	LoggerLevel    string
	LoggerJSON     bool
}

func Load(path string) (*Config, error) {
	if path != "" {
		if err := envfile.Load(path); err != nil {
			return nil, err
		}
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		return nil, fmt.Errorf("REDIS_ADDRESS is required")
	}
	ttl, err := time.ParseDuration(valueOrDefault("SESSION_TTL", "24h"))
	if err != nil {
		return nil, fmt.Errorf("parse SESSION_TTL: %w", err)
	}
	jsonLogs, err := strconv.ParseBool(valueOrDefault("LOGGER_AS_JSON", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse LOGGER_AS_JSON: %w", err)
	}
	grpcPort := valueOrDefault("GRPC_PORT", "50054")
	return &Config{
		DatabaseURL: databaseURL, RedisAddress: redisAddress, RedisPassword: os.Getenv("REDIS_PASSWORD"),
		SessionTTL: ttl, MigrationsPath: valueOrDefault("MIGRATIONS_PATH", "migrations"),
		GRPCAddress:   net.JoinHostPort(os.Getenv("GRPC_HOST"), grpcPort),
		GatewayTarget: net.JoinHostPort("127.0.0.1", grpcPort),
		HTTPAddress:   net.JoinHostPort(os.Getenv("HTTP_HOST"), valueOrDefault("HTTP_PORT", "8084")),
		LoggerLevel:   valueOrDefault("LOGGER_LEVEL", "info"), LoggerJSON: jsonLogs,
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
