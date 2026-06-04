package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

type Repository interface {
	SetOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
}
