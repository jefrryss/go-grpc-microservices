package middleware

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SessionResolver interface {
	ResolveSession(context.Context, uuid.UUID) (uuid.UUID, error)
}

type userIDKey struct{}

func UserID(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return value, ok
}

func UnaryAuth(resolver SessionResolver) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if !strings.HasSuffix(info.FullMethod, "/Whoami") {
			return handler(ctx, req)
		}
		sessionValue := metadataValue(ctx, "session-uuid")
		if sessionValue == "" {
			if request, ok := req.(interface{ GetSessionUuid() string }); ok {
				sessionValue = request.GetSessionUuid()
			}
		}
		sessionID, err := uuid.Parse(sessionValue)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, model.ErrUnauthorized.Error())
		}
		userID, err := resolver.ResolveSession(ctx, sessionID)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, model.ErrUnauthorized.Error())
		}
		return handler(context.WithValue(ctx, userIDKey{}, userID), req)
	}
}

func metadataValue(ctx context.Context, key string) string {
	values := metadata.ValueFromIncomingContext(ctx, key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
