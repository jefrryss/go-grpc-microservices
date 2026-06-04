package service

import (
	clientInventory "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory"
	clientPayment "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment"
	repository "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/service"
)

var _ service.Service = (*OrderService)(nil)

type OrderService struct {
	repo            repository.Repository
	paymentClient   clientPayment.PaymentClient
	inventoryClient clientInventory.InventoryClient
}

func NewOrderService(repo repository.Repository, payClient clientPayment.PaymentClient, inventory clientInventory.InventoryClient) *OrderService {
	return &OrderService{
		repo:            repo,
		paymentClient:   payClient,
		inventoryClient: inventory,
	}
}
