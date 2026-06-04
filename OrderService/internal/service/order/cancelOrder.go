package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (o *OrderService) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	order, err := o.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return err
	}
	if order.Status == model.OrderStatusPendingPayment {
		order.Status = model.OrderStatusCancelled
		order.UpdatedAt = time.Now()
		err := o.repo.SetOrder(ctx, order)
		if err != nil {
			return err
		}
		return nil
	}
	return model.ErrOrderCannotBeCancelled
}
