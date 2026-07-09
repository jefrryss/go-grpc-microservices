package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type SessionRepository interface {
	Create(context.Context, uuid.UUID, uuid.UUID) error
	Get(context.Context, uuid.UUID) (uuid.UUID, error)
}

type Auth struct {
	users    repository.UserRepository
	sessions SessionRepository
}

func New(users repository.UserRepository, sessions SessionRepository) *Auth {
	return &Auth{users: users, sessions: sessions}
}

func (s *Auth) Register(ctx context.Context, login, password, email string, methods []model.NotificationMethod) (uuid.UUID, error) {
	if strings.TrimSpace(login) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(email) == "" {
		return uuid.Nil, model.ErrInvalidCredentials
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}
	user := &model.User{
		ID: uuid.New(), Login: login, PasswordHash: string(hash), Email: email,
		NotificationMethods: methods,
	}
	if err := s.users.Create(ctx, user); err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (s *Auth) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	user, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return uuid.Nil, model.ErrInvalidCredentials
	}
	sessionID := uuid.New()
	if err := s.sessions.Create(ctx, sessionID, user.ID); err != nil {
		return uuid.Nil, err
	}
	return sessionID, nil
}

func (s *Auth) Whoami(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *Auth) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.users.GetByID(ctx, userID)
}

func (s *Auth) ResolveSession(ctx context.Context, sessionID uuid.UUID) (uuid.UUID, error) {
	return s.sessions.Get(ctx, sessionID)
}
