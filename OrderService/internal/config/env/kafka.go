package env

import (
	"os"
	"strings"

	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type KafkaConfig struct {
	brokers            []string
	orderPaidTopic     string
	shipAssembledTopic string
	groupID            string
}

func NewKafkaConfig() *KafkaConfig {
	var brokers []string
	if value := strings.TrimSpace(os.Getenv("KAFKA_BROKERS")); value != "" {
		brokers = strings.Split(value, ",")
	}
	return &KafkaConfig{
		brokers:            brokers,
		orderPaidTopic:     valueOrDefault("ORDER_PAID_TOPIC", events.TopicOrderPaid),
		shipAssembledTopic: valueOrDefault("SHIP_ASSEMBLED_TOPIC", events.TopicShipAssembled),
		groupID:            valueOrDefault("SHIP_ASSEMBLED_GROUP_ID", "order-service"),
	}
}

func (c *KafkaConfig) Brokers() []string            { return c.brokers }
func (c *KafkaConfig) OrderPaidTopic() string       { return c.orderPaidTopic }
func (c *KafkaConfig) ShipAssembledTopic() string   { return c.shipAssembledTopic }
func (c *KafkaConfig) ShipAssembledGroupID() string { return c.groupID }
