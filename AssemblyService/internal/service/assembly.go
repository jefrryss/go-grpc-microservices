package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type Publisher interface {
	Publish(context.Context, events.ShipAssembled) error
}

type Assembly struct {
	duration  time.Duration
	publisher Publisher
}

func NewAssembly(duration time.Duration, publisher Publisher) *Assembly {
	return &Assembly{duration: duration, publisher: publisher}
}

func (s *Assembly) Handle(ctx context.Context, event events.OrderPaid) error {
	timer := time.NewTimer(s.duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	return s.publisher.Publish(ctx, events.ShipAssembled{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(s.duration.Seconds()),
	})
}
