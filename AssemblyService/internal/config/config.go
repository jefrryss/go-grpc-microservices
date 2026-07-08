package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jefrryss/go-grpc-microservices/platform/pkg/config/envfile"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type Config struct {
	Brokers            []string
	OrderPaidTopic     string
	ShipAssembledTopic string
	GroupID            string
	BuildDuration      time.Duration
	LoggerLevel        string
	LoggerJSON         bool
}

func Load(path string) (*Config, error) {
	if path != "" {
		if err := envfile.Load(path); err != nil {
			return nil, err
		}
	}
	brokersValue := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokersValue == "" {
		return nil, fmt.Errorf("KAFKA_BROKERS is required")
	}
	buildDuration, err := time.ParseDuration(valueOrDefault("ASSEMBLY_DURATION", "10s"))
	if err != nil {
		return nil, fmt.Errorf("parse ASSEMBLY_DURATION: %w", err)
	}
	loggerJSON, err := strconv.ParseBool(valueOrDefault("LOGGER_AS_JSON", "false"))
	if err != nil {
		return nil, fmt.Errorf("parse LOGGER_AS_JSON: %w", err)
	}
	return &Config{
		Brokers:            strings.Split(brokersValue, ","),
		OrderPaidTopic:     valueOrDefault("ORDER_PAID_TOPIC", events.TopicOrderPaid),
		ShipAssembledTopic: valueOrDefault("SHIP_ASSEMBLED_TOPIC", events.TopicShipAssembled),
		GroupID:            valueOrDefault("ORDER_PAID_GROUP_ID", "assembly-service"),
		BuildDuration:      buildDuration,
		LoggerLevel:        valueOrDefault("LOGGER_LEVEL", "info"),
		LoggerJSON:         loggerJSON,
	}, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
