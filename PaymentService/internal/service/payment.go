package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
)

type Service interface {
	CreateTransaction(ctx context.Context, orderUUID, userrUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error)
}
