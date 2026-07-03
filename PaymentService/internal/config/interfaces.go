package config

type GRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	JSON() bool
}
