package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jefrryss/go-grpc-microservices/platform/pkg/config/envfile"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type Config struct {
	Brokers            []string
	OrderPaidTopic     string
	ShipAssembledTopic string
	AuthAddress        string
	TelegramToken      string
	LoggerLevel        string
	LoggerJSON         bool
}

func Load(path string) (*Config, error) {
	if path != "" {
		if err := envfile.Load(path); err != nil {
			return nil, err
		}
	}
	brokers := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokers == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	jsonLogs, err := strconv.ParseBool(valueOrDefault("LOGGER_AS_JSON", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse LOGGER_AS_JSON: %w", err)
	}
	return &Config{
		Brokers: strings.Split(brokers, ","), OrderPaidTopic: valueOrDefault("ORDER_PAID_TOPIC", events.TopicOrderPaid),
		ShipAssembledTopic: valueOrDefault("SHIP_ASSEMBLED_TOPIC", events.TopicShipAssembled),
		AuthAddress:        valueOrDefault("AUTH_SERVICE_ADDRESS", "localhost:50054"), TelegramToken: token,
		LoggerLevel: valueOrDefault("LOGGER_LEVEL", "info"), LoggerJSON: jsonLogs,
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
