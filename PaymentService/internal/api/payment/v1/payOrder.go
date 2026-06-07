package api

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/converter"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *PaymentServer) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order uuid format: %v", err)
	}
	if orderUUID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "order uuid cannot be empty")
	}

	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user uuid: %v", err)
	}
	if userUUID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "user uuid cannot be empty")
	}

	paymentMethod, err := converter.ToDomainPaymentMethod(req.GetPaymentMethod())
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid payment method: %v", err)
		}
		return nil, status.Errorf(codes.InvalidArgument, "unknown payment method code")
	}

	transactionUUID, err := p.service.CreateTransaction(ctx, orderUUID, userUUID, paymentMethod)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
