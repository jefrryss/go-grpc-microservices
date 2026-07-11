package auth

import (
	"context"
	"fmt"

	authV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc"
)

type NotificationMethod struct {
	ProviderName string
	Target       string
}

type Client struct{ client authV1.AuthServiceClient }

func New(connection grpc.ClientConnInterface) *Client {
	return &Client{client: authV1.NewAuthServiceClient(connection)}
}

func (c *Client) NotificationMethods(ctx context.Context, userUUID string) ([]NotificationMethod, error) {
	response, err := c.client.GetUser(ctx, &authV1.GetUserRequest{UserUuid: userUUID})
	if err != nil {
		return nil, fmt.Errorf("get user from AuthService: %w", err)
	}
	methods := make([]NotificationMethod, 0, len(response.GetNotificationMethods()))
	for _, method := range response.GetNotificationMethods() {
		methods = append(methods, NotificationMethod{ProviderName: method.GetProviderName(), Target: method.GetTarget()})
	}
	return methods, nil
}
