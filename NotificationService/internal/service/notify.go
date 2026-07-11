package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/jefrryss/go-grpc-microservices/NotificationService/internal/client/auth"
)

type UserClient interface {
	NotificationMethods(context.Context, string) ([]auth.NotificationMethod, error)
}

type Sender interface {
	Send(context.Context, string, string) error
}

type Notify struct {
	users    UserClient
	telegram Sender
}

func NewNotify(users UserClient, telegram Sender) *Notify {
	return &Notify{users: users, telegram: telegram}
}

func (s *Notify) User(ctx context.Context, userUUID, message string) error {
	methods, err := s.users.NotificationMethods(ctx, userUUID)
	if err != nil {
		return err
	}
	for _, method := range methods {
		if strings.EqualFold(method.ProviderName, "telegram") {
			if err := s.telegram.Send(ctx, method.Target, message); err != nil {
				return fmt.Errorf("send Telegram notification: %w", err)
			}
		}
	}
	return nil
}
