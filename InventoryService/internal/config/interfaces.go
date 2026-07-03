package config

type GRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	Database() string
	Collection() string
}

type LoggerConfig interface {
	Level() string
	JSON() bool
}
