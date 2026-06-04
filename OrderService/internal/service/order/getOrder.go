package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (o *OrderService) GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	order, err := o.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}
