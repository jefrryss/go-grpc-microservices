package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

type InventoryClient interface {
	GetListParts(ctx context.Context, parts []uuid.UUID) ([]*model.Part, error)
	GetPart(ctx context.Context, uuidPart uuid.UUID) (*model.Part, error)
}
