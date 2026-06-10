package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	inventory_v1 "InventoryService/pkg/inventory/v1"
	"OrderService/internal/domain"

	"github.com/google/uuid"
)

var (
	ErrOrderCannotBeCancelled = errors.New("order is already paid or cannot be cancelled")
	ErrOrderNotFound          = errors.New("order not found")
	ErrInvalidOrderStatus     = errors.New("invalid order status for this operation")
)

type OrderRepository interface {
	GetOrder(ctx context.Context, orderUUID uuid.UUID) (domain.Order, error)
	SetOrder(ctx context.Context, order *domain.Order) error
}

type grpcInventoryClient interface {
	GetListParts(ctx context.Context, parts []string) ([]*inventory_v1.Part, error)
	GetPart(ctx context.Context, uuidPart string) (*inventory_v1.Part, error)
}

type grpcPaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, userUUID string, paymentMethod domain.PaymentMethod) (string, error)
}

type OrderService struct {
	repo            OrderRepository
	paymentClient   grpcPaymentClient
	inventoryClient grpcInventoryClient
}

func NewOrderService(repo OrderRepository, payClient grpcPaymentClient, inventory grpcInventoryClient) *OrderService {
	return &OrderService{
		repo:            repo,
		paymentClient:   payClient,
		inventoryClient: inventory,
	}
}

func (o *OrderService) CancelOrder(ctx context.Context, orderUUID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	orderID, err := uuid.Parse(orderUUID)
	if err != nil {
		return err
	}

	order, err := o.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status == domain.OrderStatusPendingPayment {
		order.Status = domain.OrderStatusCancelled
		order.UpdatedAt = time.Now()
		err := o.repo.SetOrder(ctx, &order)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrOrderCannotBeCancelled
}

func (o *OrderService) GetOrder(ctx context.Context, orderUUID string) (*domain.Order, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	orderId, err := uuid.Parse(orderUUID)
	if err != nil {
		return nil, err
	}

	order, err := o.repo.GetOrder(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *OrderService) PayOrder(ctx context.Context, orderUUID string, paymentMethod domain.PaymentMethod) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	orderID, err := uuid.Parse(orderUUID)
	if err != nil {
		return "", err
	}

	order, err := o.repo.GetOrder(ctx, orderID)
	if err != nil {
		return "", err
	}
	if order.Status != domain.OrderStatusPendingPayment {
		return "", ErrInvalidOrderStatus
	}
	trancID, err := o.paymentClient.PayOrder(ctx, order.ID.String(), order.UserID.String(), paymentMethod)
	if err != nil {
		return "", err
	}

	trancUUID, err := uuid.Parse(trancID)
	if err != nil {
		return "", err
	}

	order.Status = domain.OrderStatusPaid
	order.TransactionID = uuid.NullUUID{UUID: trancUUID, Valid: true}
	order.PaymentMethod = paymentMethod
	order.UpdatedAt = time.Now()

	err = o.repo.SetOrder(ctx, &order)
	if err != nil {
		return "", err
	}

	return trancID, nil
}

func (o *OrderService) CreateOrder(ctx context.Context, userUUID string, partsUUIDS []string) (uuid.UUID, float64, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, 0, err
	}

	userID, err := uuid.Parse(userUUID)
	if err != nil {
		return uuid.Nil, 0, err
	}

	parts, err := o.inventoryClient.GetListParts(ctx, partsUUIDS)
	if err != nil {
		return uuid.Nil, 0, err
	}

	requestedCounts := make(map[string]int)
	for _, id := range partsUUIDS {
		requestedCounts[id]++
	}

	availableParts := make(map[string]*inventory_v1.Part)
	for _, part := range parts {
		availableParts[part.Uuid] = part
	}

	var totalPrice float64
	for reqUUIDStr, count := range requestedCounts {
		part, exists := availableParts[reqUUIDStr]
		if !exists {
			return uuid.Nil, 0, fmt.Errorf("part %s not found", reqUUIDStr)
		}
		totalPrice += float64(count) * part.Price
	}

	partsForOrder := make([]uuid.UUID, 0, len(partsUUIDS))
	for _, idStr := range partsUUIDS {
		u, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, 0, err
		}
		partsForOrder = append(partsForOrder, u)
	}

	orderUUID := uuid.New()
	order := &domain.Order{
		ID:            orderUUID,
		UserID:        userID,
		PartIDs:       partsForOrder,
		TotalPrice:    totalPrice,
		TransactionID: uuid.NullUUID{},
		Status:        domain.OrderStatusPendingPayment,
		PaymentMethod: domain.PaymentMethodUnknown,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = o.repo.SetOrder(ctx, order)
	if err != nil {
		return uuid.Nil, 0, err
	}

	return orderUUID, totalPrice, nil
}
