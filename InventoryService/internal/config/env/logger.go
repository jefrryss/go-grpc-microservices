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
	level := valueOrDefault("LOGGER_LEVEL", "info")
	jsonValue, _ := strconv.ParseBool(os.Getenv("LOGGER_AS_JSON"))
	return &LoggerConfig{level: level, json: jsonValue}
}

func (c *LoggerConfig) Level() string { return c.level }
func (c *LoggerConfig) JSON() bool    { return c.json }
