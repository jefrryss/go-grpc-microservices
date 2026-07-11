package service

import (
	"context"
	"testing"

	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/client/auth"
	"github.com/stretchr/testify/require"
)

type usersStub struct{ methods []auth.NotificationMethod }

func (u usersStub) NotificationMethods(context.Context, string) ([]auth.NotificationMethod, error) {
	return u.methods, nil
}

type senderStub struct {
	target  string
	message string
}

func (s *senderStub) Send(_ context.Context, target, message string) error {
	s.target, s.message = target, message
	return nil
}

func TestNotifyTelegramUser(t *testing.T) {
	sender := &senderStub{}
	service := NewNotify(usersStub{methods: []auth.NotificationMethod{
		{ProviderName: "email", Target: "demo@example.com"},
		{ProviderName: "telegram", Target: "123"},
	}}, sender)

	require.NoError(t, service.User(context.Background(), "user", "paid"))
	require.Equal(t, "123", sender.target)
	require.Equal(t, "paid", sender.message)
}
