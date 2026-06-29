package api

import (
	"errors"
	"testing"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapCreateOrderError(t *testing.T) {
	require.Equal(t, codes.InvalidArgument, status.Code(mapCreateOrderError(model.ErrPartNotFound)))
	require.Equal(t, codes.Unavailable, status.Code(mapCreateOrderError(status.Error(codes.Unavailable, "inventory unavailable"))))
	require.Equal(t, codes.Internal, status.Code(mapCreateOrderError(errors.New("database error"))))
}

func TestMapGetOrderError(t *testing.T) {
	require.Equal(t, codes.NotFound, status.Code(mapGetOrderError(model.ErrOrderNotFound)))
	require.Equal(t, codes.Internal, status.Code(mapGetOrderError(errors.New("database error"))))
}

func TestMapPayOrderError(t *testing.T) {
	require.Equal(t, codes.NotFound, status.Code(mapPayOrderError(model.ErrOrderNotFound)))
	require.Equal(t, codes.Aborted, status.Code(mapPayOrderError(model.ErrInvalidOrderStatus)))
	require.Equal(t, codes.Unavailable, status.Code(mapPayOrderError(status.Error(codes.Unavailable, "payment unavailable"))))
}
