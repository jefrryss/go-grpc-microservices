package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/jefrryss/go-grpc-microservices/shared/pkg/events"
)

type publisherStub struct{ event events.OrderPaid }

func (p *publisherStub) Publish(_ context.Context, event events.OrderPaid) error {
	p.event = event
	return nil
}

func (s *ServiceSuite) TestPayOrder_Success() {
	orderID := uuid.New()
	transactionID := uuid.New()
	userID := uuid.New()

	expectedOrder := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.OrderStatusPendingPayment,
	}

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(expectedOrder, nil).Once()
	s.clPaymentMock.On("PayOrder", s.ctx, orderID, userID, model.PaymentMethodCard).Return(transactionID, nil).Once()
	s.repoMock.On("SetOrder", s.ctx, expectedOrder).Return(nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.NoError(err)
	s.Equal(transactionID, res)
	s.Equal(model.OrderStatusPaid, expectedOrder.Status)
	s.Equal(model.PaymentMethodCard, expectedOrder.PaymentMethod)
	s.True(expectedOrder.TransactionID.Valid)
	s.Equal(transactionID, expectedOrder.TransactionID.UUID)
	s.WithinDuration(time.Now(), expectedOrder.UpdatedAt, time.Second)
}

func (s *ServiceSuite) TestPayOrder_PublishesEvent() {
	orderID := uuid.New()
	transactionID := uuid.New()
	userID := uuid.New()
	order := &model.Order{ID: orderID, UserID: userID, Status: model.OrderStatusPendingPayment}
	publisher := &publisherStub{}
	s.service = NewOrderService(s.repoMock, s.clPaymentMock, s.clInventoryMock, WithOrderPaidPublisher(publisher))
	s.repoMock.On("GetOrder", s.ctx, orderID).Return(order, nil).Once()
	s.clPaymentMock.On("PayOrder", s.ctx, orderID, userID, model.PaymentMethodCard).Return(transactionID, nil).Once()
	s.repoMock.On("SetOrder", s.ctx, order).Return(nil).Once()

	_, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.NoError(err)
	s.Equal(orderID.String(), publisher.event.OrderUUID)
	s.Equal(userID.String(), publisher.event.UserUUID)
	s.Equal(transactionID.String(), publisher.event.TransactionUUID)
	s.NotEmpty(publisher.event.EventUUID)
}

func (s *ServiceSuite) TestPayOrder_ContextCanceled() {
	orderID := uuid.New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := s.service.PayOrder(ctx, orderID, model.PaymentMethodCard)

	s.Equal(uuid.Nil, res)
	s.Error(err)
	s.ErrorIs(err, context.Canceled)

	s.repoMock.AssertNotCalled(s.T(), "GetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestPayOrder_GetOrderError() {
	orderID := uuid.New()
	expectedErr := errors.New("failed to get order from database")

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(nil, expectedErr).Once()

	res, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.Equal(uuid.Nil, res)
	s.Error(err)
	s.Equal(expectedErr.Error(), err.Error())
}

func (s *ServiceSuite) TestPayOrder_InvalidOrderStatus() {
	orderID := uuid.New()
	userID := uuid.New()

	expectedOrder := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.OrderStatusCancelled,
	}

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(expectedOrder, nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.Equal(uuid.Nil, res)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidOrderStatus)

	s.clPaymentMock.AssertNotCalled(s.T(), "PayOrder", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestPayOrder_PaymentClientError() {
	orderID := uuid.New()
	userID := uuid.New()
	expectedErr := errors.New("payment gateway timeout")

	expectedOrder := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.OrderStatusPendingPayment,
	}

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(expectedOrder, nil).Once()
	s.clPaymentMock.On("PayOrder", s.ctx, orderID, userID, model.PaymentMethodCard).Return(uuid.Nil, expectedErr).Once()

	res, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.Equal(uuid.Nil, res)
	s.Error(err)
	s.Equal(expectedErr.Error(), err.Error())

	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestPayOrder_SetOrderError() {
	orderID := uuid.New()
	transactionID := uuid.New()
	userID := uuid.New()
	expectedErr := errors.New("failed to write updated order status to database")

	expectedOrder := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.OrderStatusPendingPayment,
	}

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(expectedOrder, nil).Once()
	s.clPaymentMock.On("PayOrder", s.ctx, orderID, userID, model.PaymentMethodCard).Return(transactionID, nil).Once()
	s.repoMock.On("SetOrder", s.ctx, expectedOrder).Return(expectedErr).Once()

	res, err := s.service.PayOrder(s.ctx, orderID, model.PaymentMethodCard)

	s.Equal(uuid.Nil, res)
	s.Error(err)
	s.Equal(expectedErr.Error(), err.Error())
}
