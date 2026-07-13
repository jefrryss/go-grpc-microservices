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

type Observer interface {
	AssemblyCompleted(time.Duration)
}

type Assembly struct {
	duration  time.Duration
	publisher Publisher
	observer  Observer
}

func NewAssembly(duration time.Duration, publisher Publisher, observers ...Observer) *Assembly {
	assembly := &Assembly{duration: duration, publisher: publisher}
	if len(observers) > 0 {
		assembly.observer = observers[0]
	}
	return assembly
}

func (s *Assembly) Handle(ctx context.Context, event events.OrderPaid) error {
	timer := time.NewTimer(s.duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	if s.observer != nil {
		s.observer.AssemblyCompleted(s.duration)
	}
	return s.publisher.Publish(ctx, events.ShipAssembled{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(s.duration.Seconds()),
	})
}
