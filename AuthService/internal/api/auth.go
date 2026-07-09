package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/middleware"
	"github.com/jefrryss/go-grpc-microservices/AuthService/internal/model"
	authV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	Register(context.Context, string, string, string, []model.NotificationMethod) (uuid.UUID, error)
	Login(context.Context, string, string) (uuid.UUID, error)
	Whoami(context.Context, uuid.UUID) (*model.User, error)
	GetUser(context.Context, uuid.UUID) (*model.User, error)
}

type Server struct {
	authV1.UnimplementedAuthServiceServer
	service Service
}

func New(service Service) *Server { return &Server{service: service} }

func (s *Server) Register(ctx context.Context, request *authV1.RegisterRequest) (*authV1.RegisterResponse, error) {
	methods := make([]model.NotificationMethod, 0, len(request.GetNotificationMethods()))
	for _, method := range request.GetNotificationMethods() {
		methods = append(methods, model.NotificationMethod{ProviderName: method.GetProviderName(), Target: method.GetTarget()})
	}
	userID, err := s.service.Register(ctx, request.GetLogin(), request.GetPassword(), request.GetEmail(), methods)
	if err != nil {
		return nil, mapError(err)
	}
	return &authV1.RegisterResponse{UserUuid: userID.String()}, nil
}

func (s *Server) Login(ctx context.Context, request *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	sessionID, err := s.service.Login(ctx, request.GetLogin(), request.GetPassword())
	if err != nil {
		return nil, mapError(err)
	}
	return &authV1.LoginResponse{SessionUuid: sessionID.String()}, nil
}

func (s *Server) Whoami(ctx context.Context, _ *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	userID, ok := middleware.UserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, model.ErrUnauthorized.Error())
	}
	user, err := s.service.Whoami(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}
	return &authV1.WhoamiResponse{UserUuid: user.ID.String(), Login: user.Login, Email: user.Email}, nil
}

func (s *Server) GetUser(ctx context.Context, request *authV1.GetUserRequest) (*authV1.GetUserResponse, error) {
	userID, err := uuid.Parse(request.GetUserUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user UUID")
	}
	user, err := s.service.GetUser(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}
	methods := make([]*authV1.NotificationMethod, 0, len(user.NotificationMethods))
	for _, method := range user.NotificationMethods {
		methods = append(methods, &authV1.NotificationMethod{ProviderName: method.ProviderName, Target: method.Target})
	}
	return &authV1.GetUserResponse{
		UserUuid: user.ID.String(), Login: user.Login, Email: user.Email, NotificationMethods: methods,
	}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, model.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, model.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, model.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
