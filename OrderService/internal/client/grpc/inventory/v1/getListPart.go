package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/converter"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func (g *GrpcInventoryClient) GetListParts(ctx context.Context, parts []uuid.UUID) ([]*model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return []*model.Part{}, nil
	}

	strUUIDs := make([]string, 0, len(parts))
	for _, id := range parts {
		strUUIDs = append(strUUIDs, id.String())
	}

	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: strUUIDs,
		},
	}
	res, err := g.api.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}
	arrRes := make([]*model.Part, 0, len(res.Parts))
	for _, val := range res.Parts {
		domainParet, err := converter.ToDomainPart(val)
		if err != nil {
			return nil, err
		}
		arrRes = append(arrRes, domainParet)
	}
	return arrRes, nil
}
