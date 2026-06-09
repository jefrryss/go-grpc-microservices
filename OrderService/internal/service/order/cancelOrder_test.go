package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCancelOrder_CancekCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.service.CancelOrder(ctx, uuid.New())
	s.ErrorIs(err, context.Canceled)
	s.repoMock.AssertNotCalled(s.T(), "GetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCancelOrder_Success() {
	partId := uuid.New()
	order := &model.Order{
		ID:     partId,
		Status: model.OrderStatusPendingPayment,
	}
	s.repoMock.On("GetOrder", mock.Anything, mock.Anything).Return(order, nil).Once()
	s.repoMock.On("SetOrder", mock.Anything, mock.Anything).Return(nil).Once()
	err := s.service.CancelOrder(s.ctx, partId)
	expectedTime := time.Now()
	s.NoError(err)
	s.Equal(order.ID, partId)
	s.Equal(order.Status, model.OrderStatusCancelled)
	s.WithinDuration(expectedTime, order.UpdatedAt, time.Second)

}

func (s *ServiceSuite) TestCancelOrder_InvalidGetOrder() {
	errGet := errors.New("error get order")
	s.repoMock.On("GetOrder", mock.Anything, mock.Anything).Return(nil, errGet).Once()
	err := s.service.CancelOrder(s.ctx, uuid.New())
	s.ErrorIs(err, errGet)
	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCancelOrder_InvalidStatus() {
	orderId := uuid.New()
	order := &model.Order{
		ID:     orderId,
		Status: model.OrderStatusPaid,
	}

	s.repoMock.On("GetOrder", mock.Anything, mock.Anything).Return(order, nil).Once()
	err := s.service.CancelOrder(s.ctx, orderId)
	s.ErrorIs(err, model.ErrOrderCannotBeCancelled)
	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCancelOrder_errSetOrder() {
	errSetOrder := errors.New("Error set Order")
	orderId := uuid.New()
	order := &model.Order{
		ID:     orderId,
		Status: model.OrderStatusPendingPayment,
	}
	s.repoMock.On("GetOrder", mock.Anything, mock.Anything).Return(order, nil).Once()
	s.repoMock.On("SetOrder", mock.Anything, mock.Anything).Return(errSetOrder)
	err := s.service.CancelOrder(s.ctx, orderId)

	s.ErrorIs(err, errSetOrder)
}
