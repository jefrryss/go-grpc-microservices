package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
)

type UserRepository interface {
	Create(context.Context, *model.User) error
	GetByLogin(context.Context, string) (*model.User, error)
	GetByID(context.Context, uuid.UUID) (*model.User, error)
}
