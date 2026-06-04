package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/converter"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	orderId, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "order_uuid Invalid")
	}
	paymentMethod := converter.ToDomainPaymentMethod(req.GetPaymentMethod())
	trancID, err := o.service.PayOrder(ctx, orderId, paymentMethod)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to pay order: %v", err)
	}
	return &order_v1.PayOrderResponse{TransactionUuid: trancID.String()}, nil
}
