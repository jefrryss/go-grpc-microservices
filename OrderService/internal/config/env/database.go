package env

import "os"

type DatabaseConfig struct{ url string }

func NewDatabaseConfig() *DatabaseConfig { return &DatabaseConfig{url: os.Getenv("DATABASE_URL")} }
func (c *DatabaseConfig) URL() string    { return c.url }
