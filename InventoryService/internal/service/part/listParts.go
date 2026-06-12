package service

import (
	"context"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	service "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service/filter"
)

func (i *InventoryService) ListParts(ctx context.Context, filters *model.Filter) ([]*model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	partsAll, err := i.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*model.Part, 0, len(partsAll))
	filtr := service.NewInventoryFilter(filters)

	for _, part := range partsAll {
		if filtr.FilterPart(part) {
			res = append(res, part)
		}
	}
	return res, nil
}
