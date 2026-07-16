package service

import (
	"context"
	"math/rand/v2"
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
	duration := s.duration
	maxSeconds := int(s.duration / time.Second)
	if maxSeconds > 1 {
		duration = time.Duration(rand.IntN(maxSeconds)+1) * time.Second
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	if s.observer != nil {
		s.observer.AssemblyCompleted(duration)
	}
	return s.publisher.Publish(ctx, events.ShipAssembled{
		EventUUID:    uuid.NewString(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: int64(duration.Seconds()),
	})
}
