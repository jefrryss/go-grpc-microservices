package service

import (
	"context"

	clientInventory "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory"
	clientPayment "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment"
	repository "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type OrderPaidPublisher interface {
	Publish(context.Context, events.OrderPaid) error
}

var _ service.Service = (*OrderService)(nil)

type OrderService struct {
	repo            repository.Repository
	paymentClient   clientPayment.PaymentClient
	inventoryClient clientInventory.InventoryClient
	orderPaid       OrderPaidPublisher
}

func NewOrderService(repo repository.Repository, payClient clientPayment.PaymentClient, inventory clientInventory.InventoryClient, publishers ...OrderPaidPublisher) *OrderService {
	service := &OrderService{
		repo:            repo,
		paymentClient:   payClient,
		inventoryClient: inventory,
	}
	if len(publishers) > 0 {
		service.orderPaid = publishers[0]
	}
	return service
}
