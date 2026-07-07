package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (s *ServiceSuite) TestCompleteOrder_Success() {
	orderID := uuid.New()
	order := &model.Order{ID: orderID, Status: model.OrderStatusPaid}
	s.repoMock.On("GetOrder", s.ctx, orderID).Return(order, nil).Once()
	s.repoMock.On("SetOrder", s.ctx, order).Return(nil).Once()

	err := s.service.CompleteOrder(s.ctx, orderID)

	s.NoError(err)
	s.Equal(model.OrderStatusCompleted, order.Status)
	s.WithinDuration(time.Now(), order.UpdatedAt, time.Second)
}

func (s *ServiceSuite) TestCompleteOrder_InvalidStatus() {
	orderID := uuid.New()
	order := &model.Order{ID: orderID, Status: model.OrderStatusCancelled}
	s.repoMock.On("GetOrder", s.ctx, orderID).Return(order, nil).Once()

	err := s.service.CompleteOrder(s.ctx, orderID)

	s.ErrorIs(err, model.ErrInvalidOrderStatus)
}
