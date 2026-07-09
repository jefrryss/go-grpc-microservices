package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
	ttl    time.Duration
}

func New(client *redis.Client, ttl time.Duration) *Repository {
	return &Repository{client: client, ttl: ttl}
}

func (r *Repository) Create(ctx context.Context, sessionID, userID uuid.UUID) error {
	if err := r.client.Set(ctx, "session:"+sessionID.String(), userID.String(), r.ttl).Err(); err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error) {
	value, err := r.client.Get(ctx, "session:"+sessionID.String()).Result()
	if errors.Is(err, redis.Nil) {
		return uuid.Nil, model.ErrUnauthorized
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get session: %w", err)
	}
	userID, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse session user: %w", err)
	}
	return userID, nil
}
