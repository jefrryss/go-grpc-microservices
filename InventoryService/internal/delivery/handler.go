package delivery

import (
	"InventoryService/internal/repository"
	"InventoryService/internal/service"
	inventory_v1 "InventoryService/pkg/inventory/v1"
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	inventory_v1.UnimplementedInventoryServiceServer

	service *service.InventoryService
}

func NewInventoryServer(srv *service.InventoryService) *InventoryServer {
	return &InventoryServer{
		service: srv,
	}
}

func (i *InventoryServer) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid format: %v", err)
	}

	part, err := i.service.GetPart(ctx, uuid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "part not found")
		} else {
			return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
		}
	}
	return &inventory_v1.GetPartResponse{Part: part}, nil
}

func (i *InventoryServer) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := i.service.ListParts(ctx, req.Filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error: %v", err)
	}

	return &inventory_v1.ListPartsResponse{Parts: parts}, nil
}
