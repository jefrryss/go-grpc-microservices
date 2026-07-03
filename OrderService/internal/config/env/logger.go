package env

import (
	"os"
	"strconv"
)

type LoggerConfig struct {
	level string
	json  bool
}

func NewLoggerConfig() *LoggerConfig {
	jsonValue, _ := strconv.ParseBool(os.Getenv("LOGGER_AS_JSON"))
	return &LoggerConfig{level: valueOrDefault("LOGGER_LEVEL", "info"), json: jsonValue}
}

func (c *LoggerConfig) Level() string { return c.level }
func (c *LoggerConfig) JSON() bool    { return c.json }

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
