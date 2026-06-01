package service

import (
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service"
)

var _ service.Service = (*InventoryService)(nil)

type InventoryService struct {
	repo repository.Repo
}

func NewInventoryService(r repository.Repo) *InventoryService {
	return &InventoryService{
		repo: r,
	}
}
