package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/converter"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *InventoryServer) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	parsedUUID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	part, err := i.service.GetPart(ctx, parsedUUID)
	if err != nil {
		return nil, MapError(err)
	}

	partProto, err := converter.ToProtoPart(part)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal part data (ID: %s): %v", part.PartID, err)
	}

	return &inventory_v1.GetPartResponse{Part: partProto}, nil
}
