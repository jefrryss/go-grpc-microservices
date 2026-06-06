package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
)

func (i *InventoryService) GetPart(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if id == uuid.Nil {
		return nil, model.ErrNilIDPart
	}
	part, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if part == nil {
		return nil, model.ErrNotFound
	}
	return part, nil
}
