package config

import (
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/config/env"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/config/envfile"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Dependencies DependencyConfig
	Logger       LoggerConfig
	Kafka        KafkaConfig
	Tracing      TracingConfig
}

var appConfig *Config

func Load(path string) error {
	if path != "" {
		if err := envfile.Load(path); err != nil {
			return err
		}
	}

	appConfig = &Config{
		Server:       env.NewServerConfig(),
		Database:     env.NewDatabaseConfig(),
		Dependencies: env.NewDependencyConfig(),
		Logger:       env.NewLoggerConfig(),
		Kafka:        env.NewKafkaConfig(),
		Tracing:      env.NewTracingConfig(),
	}
	return nil
}

func AppConfig() *Config { return appConfig }
