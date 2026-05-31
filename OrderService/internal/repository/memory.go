package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/domain"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type OrderMemory struct {
	data map[uuid.UUID]*domain.Order
	rw   sync.RWMutex
}

func NewOrderMemory() *OrderMemory {
	return &OrderMemory{
		data: make(map[uuid.UUID]*domain.Order),
	}
}

func (o *OrderMemory) SetOrder(ctx context.Context, order *domain.Order) error {
	o.rw.Lock()
	defer o.rw.Unlock()

	o.data[order.ID] = order
	return nil
}

func (o *OrderMemory) GetOrder(ctx context.Context, orderUUID uuid.UUID) (domain.Order, error) {
	o.rw.RLock()
	defer o.rw.RUnlock()

	if order, ok := o.data[orderUUID]; ok {
		orderCopy := *order

		if order.PartIDs != nil {
			orderCopy.PartIDs = make([]uuid.UUID, len(order.PartIDs))
			copy(orderCopy.PartIDs, order.PartIDs)
		}

		return orderCopy, nil
	}
	return domain.Order{}, ErrOrderNotFound
}
