package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type OrderPaidPublisher struct{ producer platformKafka.Producer }

func NewOrderPaidPublisher(producer platformKafka.Producer) *OrderPaidPublisher {
	return &OrderPaidPublisher{producer: producer}
}

func (p *OrderPaidPublisher) Publish(ctx context.Context, event events.OrderPaid) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal OrderPaid event: %w", err)
	}
	if err := p.producer.Send(ctx, []byte(event.OrderUUID), payload); err != nil {
		return fmt.Errorf("publish OrderPaid event: %w", err)
	}
	return nil
}
