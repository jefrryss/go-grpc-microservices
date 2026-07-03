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
