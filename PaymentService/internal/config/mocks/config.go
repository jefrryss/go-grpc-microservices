package mocks

type GRPCConfig struct {
	Value string
}

func (c GRPCConfig) Address() string { return c.Value }

type LoggerConfig struct {
	LevelValue string
	JSONValue  bool
}

func (c LoggerConfig) Level() string { return c.LevelValue }
func (c LoggerConfig) JSON() bool    { return c.JSONValue }

type TracingConfig struct{ EndpointValue string }

func (c TracingConfig) Endpoint() string { return c.EndpointValue }
