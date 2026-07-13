package config

type ServerConfig interface {
	GRPCAddress() string
	GatewayEndpoint() string
	HTTPAddress() string
}

type DatabaseConfig interface {
	URL() string
}

type DependencyConfig interface {
	InventoryAddress() string
	PaymentAddress() string
}

type LoggerConfig interface {
	Level() string
	JSON() bool
}

type KafkaConfig interface {
	Brokers() []string
	OrderPaidTopic() string
	ShipAssembledTopic() string
	ShipAssembledGroupID() string
}

type TracingConfig interface {
	Endpoint() string
}
