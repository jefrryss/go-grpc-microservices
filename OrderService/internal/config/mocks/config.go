package mocks

type ServerConfig struct {
	GRPCValue    string
	GatewayValue string
	HTTPValue    string
}

func (c ServerConfig) GRPCAddress() string     { return c.GRPCValue }
func (c ServerConfig) GatewayEndpoint() string { return c.GatewayValue }
func (c ServerConfig) HTTPAddress() string     { return c.HTTPValue }

type DatabaseConfig struct{ Value string }

func (c DatabaseConfig) URL() string { return c.Value }

type DependencyConfig struct {
	InventoryValue string
	PaymentValue   string
}

func (c DependencyConfig) InventoryAddress() string { return c.InventoryValue }
func (c DependencyConfig) PaymentAddress() string   { return c.PaymentValue }

type LoggerConfig struct {
	LevelValue string
	JSONValue  bool
}

func (c LoggerConfig) Level() string { return c.LevelValue }
func (c LoggerConfig) JSON() bool    { return c.JSONValue }

type KafkaConfig struct {
	BrokerValues       []string
	OrderPaidValue     string
	ShipAssembledValue string
	GroupValue         string
}

func (c KafkaConfig) Brokers() []string            { return c.BrokerValues }
func (c KafkaConfig) OrderPaidTopic() string       { return c.OrderPaidValue }
func (c KafkaConfig) ShipAssembledTopic() string   { return c.ShipAssembledValue }
func (c KafkaConfig) ShipAssembledGroupID() string { return c.GroupValue }

type TracingConfig struct{ EndpointValue string }

func (c TracingConfig) Endpoint() string { return c.EndpointValue }
