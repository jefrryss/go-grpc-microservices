package service

import "github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service"

var _ service.Service = (*PaymentService)(nil)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}
