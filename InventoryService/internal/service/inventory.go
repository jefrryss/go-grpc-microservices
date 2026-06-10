package service

import (
	"InventoryService/internal/repository"
	inventory_v1 "InventoryService/pkg/inventory/v1"
	"context"

	"github.com/google/uuid"
)

type InventoryService struct {
	repo repository.Repo
}

func NewInventoryService(r repository.Repo) *InventoryService {
	return &InventoryService{
		repo: r,
	}
}

func (i *InventoryService) GetPart(ctx context.Context, id uuid.UUID) (*inventory_v1.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	part, err := i.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return part, nil
}

func (i *InventoryService) GetAll(ctx context.Context) ([]*inventory_v1.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	partsAll := i.repo.GetAll(ctx)
	return partsAll, nil
}

func (i *InventoryService) ListParts(ctx context.Context, filters *inventory_v1.PartsFilter) ([]*inventory_v1.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	partsAll := i.repo.GetAll(ctx)
	res := make([]*inventory_v1.Part, 0, len(partsAll))
	filtr := NewInventoryFilter(filters)

	for _, part := range partsAll {
		if filtr.FilterPart(part) {
			res = append(res, part)
		}
	}
	return res, nil
}
