package client

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/converter"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func (g *GrpcInventoryClient) GetPart(ctx context.Context, uuidPart uuid.UUID) (*model.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req := &inventory_v1.GetPartRequest{
		Uuid: uuidPart.String(),
	}
	res, err := g.api.GetPart(ctx, req)
	if err != nil {
		return nil, err
	}
	domainPArt, err := converter.ToDomainPart(res.Part)
	return domainPArt, nil
}
