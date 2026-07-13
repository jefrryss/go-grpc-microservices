package env

import "os"

type TracingConfig struct{ endpoint string }

func NewTracingConfig() *TracingConfig {
	return &TracingConfig{endpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")}
}

func (c *TracingConfig) Endpoint() string { return c.endpoint }
