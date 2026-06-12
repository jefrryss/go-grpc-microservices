package repository

import (
	"context"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/converter"
)

func (o *OrderMemory) SetOrder(ctx context.Context, order *model.Order) error {
	o.rw.Lock()
	defer o.rw.Unlock()

	repoOrder := converter.ToRepoOrder(order)

	if repoOrder != nil {
		o.data[repoOrder.ID] = repoOrder
	}
	return nil
}
