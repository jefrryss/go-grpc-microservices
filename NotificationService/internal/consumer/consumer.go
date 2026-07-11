package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type Notifier interface {
	User(context.Context, string, string) error
}

type Consumer struct {
	consumer platformKafka.Consumer
	handler  platformKafka.Handler
}

func NewOrderPaid(consumer platformKafka.Consumer, notifier Notifier) *Consumer {
	return &Consumer{consumer: consumer, handler: func(ctx context.Context, message platformKafka.Message) error {
		var event events.OrderPaid
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("decode OrderPaid event: %w", err)
		}
		return notifier.User(ctx, event.UserUUID, "Заказ "+event.OrderUUID+" оплачен и передан в сборку.")
	}}
}

func NewShipAssembled(consumer platformKafka.Consumer, notifier Notifier) *Consumer {
	return &Consumer{consumer: consumer, handler: func(ctx context.Context, message platformKafka.Message) error {
		var event events.ShipAssembled
		if err := json.Unmarshal(message.Value, &event); err != nil {
			return fmt.Errorf("decode ShipAssembled event: %w", err)
		}
		return notifier.User(ctx, event.UserUUID, "Заказ "+event.OrderUUID+" собран.")
	}}
}

func (c *Consumer) Run(ctx context.Context) error { return c.consumer.Consume(ctx, c.handler) }
