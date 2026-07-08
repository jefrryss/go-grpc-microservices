package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type OrderPaidHandler interface {
	Handle(context.Context, events.OrderPaid) error
}

type Consumer struct {
	consumer platformKafka.Consumer
	handler  OrderPaidHandler
}

func NewConsumer(consumer platformKafka.Consumer, handler OrderPaidHandler) *Consumer {
	return &Consumer{consumer: consumer, handler: handler}
}

func (c *Consumer) Run(ctx context.Context) error {
	return c.consumer.Consume(ctx, func(ctx context.Context, message platformKafka.Message) error {
		var event events.OrderPaid
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("decode OrderPaid event: %w", err)
		}
		return c.handler.Handle(ctx, event)
	})
}

type Publisher struct{ producer platformKafka.Producer }

func NewPublisher(producer platformKafka.Producer) *Publisher {
	return &Publisher{producer: producer}
}

func (p *Publisher) Publish(ctx context.Context, event events.ShipAssembled) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal ShipAssembled event: %w", err)
	}
	if err := p.producer.Send(ctx, []byte(event.OrderUUID), payload); err != nil {
		return fmt.Errorf("publish ShipAssembled event: %w", err)
	}
	return nil
}
