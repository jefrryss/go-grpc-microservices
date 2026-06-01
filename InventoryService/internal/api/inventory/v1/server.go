package api

import (
	"errors"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryServer struct {
	inventory_v1.UnimplementedInventoryServiceServer

	service service.Service
}

func NewInventoryServer(srv service.Service) *InventoryServer {
	return &InventoryServer{
		service: srv,
	}
}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, model.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, model.ErrInvalidCategory) || errors.Is(err, model.ErrInvalidFilter) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Errorf(codes.Internal, "internal server error: %v", err)
}
