package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (o *OrderService) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}

	order, err := o.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, err
	}
	if order.Status != model.OrderStatusPendingPayment {
		return uuid.Nil, model.ErrInvalidOrderStatus
	}
	trancID, err := o.paymentClient.PayOrder(ctx, orderUUID, order.UserID, paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	order.Status = model.OrderStatusPaid
	order.TransactionID = uuid.NullUUID{UUID: trancID, Valid: true}
	order.PaymentMethod = paymentMethod
	order.UpdatedAt = time.Now()

	err = o.repo.SetOrder(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}

	return trancID, nil
}
