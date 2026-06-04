package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, orderID, userID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
}
