package delivery

import (
	"OrderService/internal/service"
	order_v1 "OrderService/pkg/proto/order/v1"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	order_v1.UnimplementedOrderServiceServer
	service *service.OrderService
}

func NewOrderServer(svc *service.OrderService) *OrderServer {
	return &OrderServer{
		service: svc,
	}
}

func (o *OrderServer) CancelOrder(ctx context.Context, req *order_v1.CancelOrderRequest) (*order_v1.CancelOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	err := o.service.CancelOrder(ctx, req.GetOrderUuid())
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
		}
		if errors.Is(err, service.ErrOrderCannotBeCancelled) {
			return nil, status.Errorf(codes.Aborted, "conflict: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return &order_v1.CancelOrderResponse{}, nil
}

func (o *OrderServer) GetOrderByUUID(ctx context.Context, req *order_v1.GetOrderByUUIDRequest) (*order_v1.GetOrderByUUIDResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	order, err := o.service.GetOrder(ctx, req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &order_v1.GetOrderByUUIDResponse{
		Order: ToProtoOrder(order),
	}, nil
}

func (o *OrderServer) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest) (*order_v1.PayOrderResponse, error) {
	if err := ctx.Err(); err != nil {
		return nil, status.Error(codes.Canceled, "request was canceled by the client")
	}

	if req.GetOrderUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}

	if req.PaymentMethod == order_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "payment method cant be unknown")
	}
	trancID, err := o.service.PayOrder(ctx, req.GetOrderUuid(), ToDomainPaymentMethod(req.GetPaymentMethod()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to pay order: %v", err)
	}
	return &order_v1.PayOrderResponse{TransactionUuid: trancID}, nil
}

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
	uuid, totalPrice, err := o.service.CreateOrder(ctx, request.GetUserUuid(), request.GetPartUuids())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create order: %v", err)
	}

	return &order_v1.CreateOrderResponse{OrderUuid: uuid.String(), TotalPrice: totalPrice}, nil
}
