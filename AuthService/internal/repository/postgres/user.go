package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
)

type UserRepository struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	methods, err := json.Marshal(user.NotificationMethods)
	if err != nil {
		return fmt.Errorf("marshal notification methods: %w", err)
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO users (id, login, password_hash, email, notification_methods)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Login, user.PasswordHash, user.Email, methods)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return model.ErrUserAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	return r.get(ctx, `SELECT id, login, password_hash, email, notification_methods FROM users WHERE login = $1`, login)
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return r.get(ctx, `SELECT id, login, password_hash, email, notification_methods FROM users WHERE id = $1`, id)
}

func (r *UserRepository) get(ctx context.Context, query string, value any) (*model.User, error) {
	var user model.User
	var methods []byte
	err := r.db.QueryRow(ctx, query, value).Scan(
		&user.ID, &user.Login, &user.PasswordHash, &user.Email, &methods,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if err := json.Unmarshal(methods, &user.NotificationMethods); err != nil {
		return nil, fmt.Errorf("decode notification methods: %w", err)
	}
	return &user, nil
}
