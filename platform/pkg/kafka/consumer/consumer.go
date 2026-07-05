package consumer

import (
	"context"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ platformKafka.Consumer = (*Consumer)(nil)

type Middleware func(platformKafka.Handler) platformKafka.Handler

type Consumer struct {
	group       sarama.ConsumerGroup
	topics      []string
	logger      *logger.Logger
	middlewares []Middleware
}

func New(brokers []string, groupID string, topics []string, log *logger.Logger, middlewares ...Middleware) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("create Kafka consumer group: %w", err)
	}
	return &Consumer{group: group, topics: topics, logger: log, middlewares: middlewares}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler platformKafka.Handler) error {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	groupHandler := &groupHandler{handler: handler, logger: c.logger}

	for ctx.Err() == nil {
		if err := c.group.Consume(ctx, c.topics, groupHandler); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) || errors.Is(err, context.Canceled) {
				return nil
			}
			return fmt.Errorf("consume Kafka messages: %w", err)
		}
	}
	return nil
}

func (c *Consumer) Close() error { return c.group.Close() }

type groupHandler struct {
	handler platformKafka.Handler
	logger  *logger.Logger
}

func (h *groupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *groupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *groupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		wrapped := platformKafka.Message{
			Key: message.Key, Value: message.Value, Topic: message.Topic,
			Partition: message.Partition, Offset: message.Offset,
		}
		if err := h.handler(session.Context(), wrapped); err != nil {
			h.logger.Error(session.Context(), "Kafka handler failed", zap.Error(err))
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}
