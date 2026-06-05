package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Part, error)
	GetAll(ctx context.Context) ([]*model.Part, error)
}
