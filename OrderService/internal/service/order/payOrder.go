package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/tracing"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

func (o *OrderService) PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	_, span := tracing.Start(ctx, "order-service", "PayOrder")
	defer span.End()
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}

	order, err := o.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return uuid.Nil, err
	}
	if o.observer != nil {
		o.observer.OrderPaid(order.TotalPrice)
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
	if o.orderPaid != nil {
		if err := o.orderPaid.Publish(ctx, events.OrderPaid{
			EventUUID:       uuid.NewString(),
			OrderUUID:       order.ID.String(),
			UserUUID:        order.UserID.String(),
			PaymentMethod:   string(paymentMethod),
			TransactionUUID: trancID.String(),
		}); err != nil {
			return uuid.Nil, err
		}
	}

	return trancID, nil
}
