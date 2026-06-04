package converter

import (
	"github.com/google/uuid"
	domain "github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	repo "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/model"
)

func ToRepoOrder(order *domain.Order) *repo.OrderRepo {
	if order == nil {
		return nil
	}

	var partIDs []uuid.UUID
	if order.PartIDs != nil {
		partIDs = make([]uuid.UUID, len(order.PartIDs))
		copy(partIDs, order.PartIDs)
	}

	return &repo.OrderRepo{
		ID:            order.ID,
		UserID:        order.UserID,
		PartIDs:       partIDs,
		TotalPrice:    order.TotalPrice,
		TransactionID: order.TransactionID,
		PaymentMethod: string(order.PaymentMethod),
		Status:        string(order.Status),
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}
}

func ToDomainOrder(orderRepo *repo.OrderRepo) *domain.Order {
	if orderRepo == nil {
		return nil
	}

	var partIDs []uuid.UUID
	if orderRepo.PartIDs != nil {
		partIDs = make([]uuid.UUID, len(orderRepo.PartIDs))
		copy(partIDs, orderRepo.PartIDs)
	}

	return &domain.Order{
		ID:            orderRepo.ID,
		UserID:        orderRepo.UserID,
		PartIDs:       partIDs,
		TotalPrice:    orderRepo.TotalPrice,
		TransactionID: orderRepo.TransactionID,
		PaymentMethod: domain.PaymentMethod(orderRepo.PaymentMethod),
		Status:        domain.OrderStatus(orderRepo.Status),
		CreatedAt:     orderRepo.CreatedAt,
		UpdatedAt:     orderRepo.UpdatedAt,
	}
}
