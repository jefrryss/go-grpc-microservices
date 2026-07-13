package config

import (
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/config/env"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/config/envfile"
)

type Config struct {
	GRPC    GRPCConfig
	Logger  LoggerConfig
	Tracing TracingConfig
}

var appConfig *Config

func Load(path string) error {
	if path != "" {
		if err := envfile.Load(path); err != nil {
			return err
		}
	}

	appConfig = &Config{
		GRPC:    env.NewGRPCConfig(),
		Logger:  env.NewLoggerConfig(),
		Tracing: env.NewTracingConfig(),
	}
	return nil
}

func AppConfig() *Config {
	return appConfig
}
