package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (o *OrderService) CreateOrder(ctx context.Context, userUUID uuid.UUID, partsUUIDS []uuid.UUID) (uuid.UUID, float64, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, 0, err
	}

	parts, err := o.inventoryClient.GetListParts(ctx, partsUUIDS)
	if err != nil {
		return uuid.Nil, 0, err
	}

	requestedCounts := make(map[uuid.UUID]int)
	for _, id := range partsUUIDS {
		requestedCounts[id]++
	}

	availableParts := make(map[uuid.UUID]*model.Part)
	for _, part := range parts {
		availableParts[part.ID] = part
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
	for _, id := range partsUUIDS {
		partsForOrder = append(partsForOrder, id)
	}

	orderUUID := uuid.New()
	order := &model.Order{
		ID:            orderUUID,
		UserID:        userUUID,
		PartIDs:       partsForOrder,
		TotalPrice:    totalPrice,
		TransactionID: uuid.NullUUID{},
		Status:        model.OrderStatusPendingPayment,
		PaymentMethod: model.PaymentMethodUnknown,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = o.repo.SetOrder(ctx, order)
	if err != nil {
		return uuid.Nil, 0, err
	}

	return orderUUID, totalPrice, nil
}
