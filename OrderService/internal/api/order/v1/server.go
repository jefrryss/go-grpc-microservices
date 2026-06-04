package api

import (
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/service"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
)

type OrderServer struct {
	order_v1.UnimplementedOrderServiceServer
	service service.Service
}

func NewOrderServer(svc service.Service) *OrderServer {
	return &OrderServer{
		service: svc,
	}
}
