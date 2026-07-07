package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type OrderCompleter interface {
	CompleteOrder(context.Context, uuid.UUID) error
}

type ShipAssembledConsumer struct {
	consumer platformKafka.Consumer
	service  OrderCompleter
}

func NewShipAssembledConsumer(consumer platformKafka.Consumer, service OrderCompleter) *ShipAssembledConsumer {
	return &ShipAssembledConsumer{consumer: consumer, service: service}
}

func (c *ShipAssembledConsumer) Run(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handle)
}

func (c *ShipAssembledConsumer) handle(ctx context.Context, message platformKafka.Message) error {
	var event events.ShipAssembled
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("decode ShipAssembled event: %w", err)
	}
	orderID, err := uuid.Parse(event.OrderUUID)
	if err != nil {
		return fmt.Errorf("parse order UUID: %w", err)
	}
	if err := c.service.CompleteOrder(ctx, orderID); err != nil {
		return fmt.Errorf("complete order: %w", err)
	}
	return nil
}
