package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
)

var _ repository.Repo = (*MemoryRepo)(nil)

type MemoryRepo struct {
	data map[uuid.UUID]*repoModel.PartRepo
	rw   sync.RWMutex
}

func NewMemoryRepo() *MemoryRepo {
	m := &MemoryRepo{
		data: make(map[uuid.UUID]*repoModel.PartRepo),
	}
	return m
}
