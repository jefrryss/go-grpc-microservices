package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/converter"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *PaymentServer) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order uuid: %v", err)
	}

	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user uuid: %v", err)
	}

	paymentMethod, err := converter.ToDomainPaymentMethod(req.GetPaymentMethod())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "payment method error: %v", err)
	}

	transactionUUID, err := p.service.CreateTransaction(ctx, orderUUID, userUUID, paymentMethod)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
