package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (o *OrderService) CompleteOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := o.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return err
	}
	if order.Status != model.OrderStatusPaid {
		return model.ErrInvalidOrderStatus
	}
	order.Status = model.OrderStatusAssembled
	order.UpdatedAt = time.Now()
	return o.repo.SetOrder(ctx, order)
}
