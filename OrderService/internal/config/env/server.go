package env

import (
	"net"
	"os"
)

type ServerConfig struct {
	grpcHost string
	grpcPort string
	httpHost string
	httpPort string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		grpcHost: os.Getenv("GRPC_HOST"),
		grpcPort: valueOrDefault("GRPC_PORT", "50051"),
		httpHost: os.Getenv("HTTP_HOST"),
		httpPort: valueOrDefault("HTTP_PORT", "8080"),
	}
}

func (c *ServerConfig) GRPCAddress() string { return net.JoinHostPort(c.grpcHost, c.grpcPort) }
func (c *ServerConfig) HTTPAddress() string { return net.JoinHostPort(c.httpHost, c.httpPort) }
func (c *ServerConfig) GatewayEndpoint() string {
	return net.JoinHostPort("127.0.0.1", c.grpcPort)
}
