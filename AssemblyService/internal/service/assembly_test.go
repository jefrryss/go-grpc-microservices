package service

import (
	"context"
	"testing"
	"time"

	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
	"github.com/stretchr/testify/require"
)

type publisherStub struct{ event events.ShipAssembled }

func (p *publisherStub) Publish(_ context.Context, event events.ShipAssembled) error {
	p.event = event
	return nil
}

func TestAssemblyHandle(t *testing.T) {
	publisher := &publisherStub{}
	assembly := NewAssembly(time.Millisecond, publisher)

	err := assembly.Handle(context.Background(), events.OrderPaid{OrderUUID: "order", UserUUID: "user"})

	require.NoError(t, err)
	require.Equal(t, "order", publisher.event.OrderUUID)
	require.Equal(t, "user", publisher.event.UserUUID)
	require.NotEmpty(t, publisher.event.EventUUID)
}
