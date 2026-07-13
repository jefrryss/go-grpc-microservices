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

type Observer interface {
	OrderCreated()
	OrderPaid(float64)
}

type Option func(*OrderService)

func WithOrderPaidPublisher(publisher OrderPaidPublisher) Option {
	return func(service *OrderService) { service.orderPaid = publisher }
}

func WithObserver(observer Observer) Option {
	return func(service *OrderService) { service.observer = observer }
}

var _ service.Service = (*OrderService)(nil)

type OrderService struct {
	repo            repository.Repository
	paymentClient   clientPayment.PaymentClient
	inventoryClient clientInventory.InventoryClient
	orderPaid       OrderPaidPublisher
	observer        Observer
}

func NewOrderService(repo repository.Repository, payClient clientPayment.PaymentClient, inventory clientInventory.InventoryClient, options ...Option) *OrderService {
	service := &OrderService{
		repo:            repo,
		paymentClient:   payClient,
		inventoryClient: inventory,
	}
	for _, option := range options {
		option(service)
	}
	return service
}
