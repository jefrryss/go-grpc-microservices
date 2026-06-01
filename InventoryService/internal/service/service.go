package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
)

type Service interface {
	ListParts(ctx context.Context, filters *model.Filter) ([]*model.Part, error)
	GetPart(ctx context.Context, id uuid.UUID) (*model.Part, error)
}
