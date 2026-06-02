package api

import (
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
)

type PaymentServer struct {
	service service.Service
	payment_v1.UnimplementedPaymentServiceServer
}

func NewPaymentServer(service service.Service) *PaymentServer {
	return &PaymentServer{
		service: service,
	}
}
