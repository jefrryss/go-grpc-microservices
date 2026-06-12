package repository

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository"
	repoModel "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/model"
)

var _ repository.Repository = (*OrderMemory)(nil)

type OrderMemory struct {
	data map[uuid.UUID]*repoModel.OrderRepo
	rw   sync.RWMutex
}

func NewOrderMemory() *OrderMemory {
	return &OrderMemory{
		data: make(map[uuid.UUID]*repoModel.OrderRepo),
	}
}
