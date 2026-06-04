package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

type Service interface {
	CancelOrder(ctx context.Context, orderUUID uuid.UUID) error
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
	PayOrder(ctx context.Context, orderUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
	CreateOrder(ctx context.Context, userUUID uuid.UUID, partsUUIDS []uuid.UUID) (uuid.UUID, float64, error)
}
