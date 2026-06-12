package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/converter"
)

func (o *OrderMemory) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	o.rw.RLock()
	defer o.rw.RUnlock()

	if repoOrder, ok := o.data[orderUUID]; ok {
		domainOrder := converter.ToDomainOrder(repoOrder)
		if domainOrder != nil {
			return domainOrder, nil
		}
	}

	return &model.Order{}, model.ErrOrderNotFound
}
