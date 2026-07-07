package messaging

import (
	"context"
	"testing"

	"github.com/google/uuid"
	platformKafka "github.com/jefrryss/go-grpc-microservices/platform/pkg/kafka"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
	"github.com/stretchr/testify/require"
)

type producerStub struct{ value []byte }

func (p *producerStub) Send(_ context.Context, _, value []byte) error { p.value = value; return nil }
func (p *producerStub) Close() error                                  { return nil }

type consumerStub struct{ message platformKafka.Message }

func (c *consumerStub) Consume(ctx context.Context, handler platformKafka.Handler) error {
	return handler(ctx, c.message)
}
func (c *consumerStub) Close() error { return nil }

type completerStub struct{ orderID uuid.UUID }

func (c *completerStub) CompleteOrder(_ context.Context, orderID uuid.UUID) error {
	c.orderID = orderID
	return nil
}

func TestOrderPaidPublisher(t *testing.T) {
	producer := &producerStub{}
	publisher := NewOrderPaidPublisher(producer)
	event := events.OrderPaid{EventUUID: uuid.NewString(), OrderUUID: uuid.NewString()}

	require.NoError(t, publisher.Publish(context.Background(), event))
	require.Contains(t, string(producer.value), event.OrderUUID)
}

func TestShipAssembledConsumer(t *testing.T) {
	orderID := uuid.New()
	consumer := &consumerStub{message: platformKafka.Message{Value: []byte(`{"order_uuid":"` + orderID.String() + `"}`)}}
	completer := &completerStub{}

	require.NoError(t, NewShipAssembledConsumer(consumer, completer).Run(context.Background()))
	require.Equal(t, orderID, completer.orderID)
}
