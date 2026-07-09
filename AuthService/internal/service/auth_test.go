package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type usersStub struct {
	created *model.User
	user    *model.User
}

func (u *usersStub) Create(_ context.Context, user *model.User) error { u.created = user; return nil }
func (u *usersStub) GetByLogin(context.Context, string) (*model.User, error) {
	return u.user, nil
}
func (u *usersStub) GetByID(context.Context, uuid.UUID) (*model.User, error) { return u.user, nil }

type sessionsStub struct {
	sessionID uuid.UUID
	userID    uuid.UUID
}

func (s *sessionsStub) Create(_ context.Context, sessionID, userID uuid.UUID) error {
	s.sessionID, s.userID = sessionID, userID
	return nil
}
func (s *sessionsStub) Get(context.Context, uuid.UUID) (uuid.UUID, error) { return s.userID, nil }

func TestRegisterHashesPassword(t *testing.T) {
	users := &usersStub{}
	auth := New(users, &sessionsStub{})

	userID, err := auth.Register(context.Background(), "demo", "secret", "demo@example.com", nil)

	require.NoError(t, err)
	require.Equal(t, userID, users.created.ID)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(users.created.PasswordHash), []byte("secret")))
}

func TestLoginCreatesSession(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	require.NoError(t, err)
	user := &model.User{ID: uuid.New(), Login: "demo", PasswordHash: string(hash)}
	sessions := &sessionsStub{}
	auth := New(&usersStub{user: user}, sessions)

	sessionID, err := auth.Login(context.Background(), "demo", "secret")

	require.NoError(t, err)
	require.Equal(t, sessionID, sessions.sessionID)
	require.Equal(t, user.ID, sessions.userID)
}

func TestLoginRejectsPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	require.NoError(t, err)
	auth := New(&usersStub{user: &model.User{PasswordHash: string(hash)}}, &sessionsStub{})

	_, err = auth.Login(context.Background(), "demo", "wrong")

	require.ErrorIs(t, err, model.ErrInvalidCredentials)
}
