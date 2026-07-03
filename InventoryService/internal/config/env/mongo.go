package env

import "os"

type MongoConfig struct {
	uri        string
	database   string
	collection string
}

func NewMongoConfig() *MongoConfig {
	return &MongoConfig{
		uri:        os.Getenv("MONGO_URI"),
		database:   valueOrDefault("MONGO_DATABASE", "inventory"),
		collection: valueOrDefault("MONGO_COLLECTION", "parts"),
	}
}

func (c *MongoConfig) URI() string        { return c.uri }
func (c *MongoConfig) Database() string   { return c.database }
func (c *MongoConfig) Collection() string { return c.collection }

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
