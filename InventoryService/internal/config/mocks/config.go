package mocks

type GRPCConfig struct{ Value string }

func (c GRPCConfig) Address() string { return c.Value }

type MongoConfig struct {
	URIValue        string
	DatabaseValue   string
	CollectionValue string
}

func (c MongoConfig) URI() string        { return c.URIValue }
func (c MongoConfig) Database() string   { return c.DatabaseValue }
func (c MongoConfig) Collection() string { return c.CollectionValue }

type LoggerConfig struct {
	LevelValue string
	JSONValue  bool
}

func (c LoggerConfig) Level() string { return c.LevelValue }
func (c LoggerConfig) JSON() bool    { return c.JSONValue }
