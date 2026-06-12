package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *OrderServer) CancelOrder(ctx context.Context, req *order_v1.CancelOrderRequest) (*order_v1.CancelOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	orderId, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "order_uuid invalid")
	}
	err = o.service.CancelOrder(ctx, orderId)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
		}
		if errors.Is(err, model.ErrOrderCannotBeCancelled) {
			return nil, status.Errorf(codes.Aborted, "conflict: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &order_v1.CancelOrderResponse{}, nil
}
