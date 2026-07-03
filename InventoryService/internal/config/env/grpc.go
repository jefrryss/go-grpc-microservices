package env

import (
	"net"
	"os"
)

type GRPCConfig struct {
	host string
	port string
}

func NewGRPCConfig() *GRPCConfig {
	host := os.Getenv("GRPC_HOST")
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}
	return &GRPCConfig{host: host, port: port}
}

func (c *GRPCConfig) Address() string { return net.JoinHostPort(c.host, c.port) }
