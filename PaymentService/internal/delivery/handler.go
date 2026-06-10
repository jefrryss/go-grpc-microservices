package delivery

import (
	"PaymentService/internal/service"
	payment_v1 "PaymentService/pkg/payment/v1"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServer struct {
	service *service.PaymentService
	payment_v1.UnimplementedPaymentServiceServer
}

func NewPaymentServer(service *service.PaymentService) *PaymentServer {
	return &PaymentServer{
		service: service,
	}
}

func (p *PaymentServer) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	uuidTranc, err := p.service.CreateTrunc(ctx, req.OrderUuid, req.UserUuid, req.PaymentMethod)

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create transaction: %v", err)
	}
	return &payment_v1.PayOrderResponse{TransactionUuid: uuidTranc}, nil
}
