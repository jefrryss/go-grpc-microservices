package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/converter"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *OrderServer) GetOrderByUUID(ctx context.Context, req *order_v1.GetOrderByUUIDRequest) (*order_v1.GetOrderByUUIDResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}
	orderId, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "order_uuid Invalid")
	}
	order, err := o.service.GetOrder(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &order_v1.GetOrderByUUIDResponse{
		Order: converter.ToProtoOrder(order),
	}, nil
}
