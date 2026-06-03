package api

import (
	"context"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/converter"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *InventoryServer) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	modelFilter, err := converter.ConvertToDomainFilter(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid filter parameters: %v", err)
	}

	parts, err := i.service.ListParts(ctx, modelFilter)
	if err != nil {
		return nil, MapError(err)
	}

	partProtoArr := make([]*inventory_v1.Part, 0, len(parts))
	for _, val := range parts {
		protoPart, err := converter.ToProtoPart(val)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to marshal part data (ID: %s): %v", val.PartID, err)
		}
		partProtoArr = append(partProtoArr, protoPart)
	}

	return &inventory_v1.ListPartsResponse{Parts: partProtoArr}, nil
}
