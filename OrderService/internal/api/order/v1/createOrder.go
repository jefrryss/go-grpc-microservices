package api

import (
	"context"

	"github.com/google/uuid"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (o *OrderServer) CreateOrder(ctx context.Context, request *order_v1.CreateOrderRequest) (*order_v1.CreateOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if request.GetUserUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	if len(request.GetPartUuids()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "part_uuids list cannot be empty")
	}
	useID, err := uuid.Parse(request.GetUserUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "user_uuid Invalid")
	}
	partUuids := make([]uuid.UUID, 0, len(request.GetPartUuids()))
	for _, uuidStr := range request.GetPartUuids() {
		uuid, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid part id: %s", uuidStr)
		}
		partUuids = append(partUuids, uuid)
	}
	uuid, totalPrice, err := o.service.CreateOrder(ctx, useID, partUuids)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create order: %v", err)
	}

	return &order_v1.CreateOrderResponse{OrderUuid: uuid.String(), TotalPrice: totalPrice}, nil
}
