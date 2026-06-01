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
	part, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return part, nil
}
