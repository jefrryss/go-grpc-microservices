package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
	"go.uber.org/zap"
)

type loggerStub struct{ message string }

func (l *loggerStub) Info(_ context.Context, message string, _ ...zap.Field) { l.message = message }

func (s *PaymentServiceSuite) TestCreateTranc_Success() {
	orderUuid := uuid.New()
	userUuid := uuid.New()
	paymentMethod := model.PaymentMethodCard

	trancUUID, err := s.service.CreateTransaction(s.ctx, orderUuid, userUuid, paymentMethod)
	s.NoError(err)
	s.NotEqual(trancUUID, uuid.Nil)
}

func (s *PaymentServiceSuite) TestCreateTranc_InvalidUserUUID() {
	orderUuid := uuid.New()
	userUuid := uuid.Nil
	paymentMethod := model.PaymentMethodCard

	trancUUID, err := s.service.CreateTransaction(s.ctx, orderUuid, userUuid, paymentMethod)
	s.ErrorIs(err, model.ErrEmptyUserUUID)
	s.Equal(trancUUID, uuid.Nil)
}

func (s *PaymentServiceSuite) TestCreateTranc_InvalidOrderUUID() {
	orderUuid := uuid.Nil
	userUuid := uuid.New()
	paymentMethod := model.PaymentMethodCard

	trancUUID, err := s.service.CreateTransaction(s.ctx, orderUuid, userUuid, paymentMethod)
	s.ErrorIs(err, model.ErrEmptyOrderUUID)
	s.Equal(trancUUID, uuid.Nil)
}

func (s *PaymentServiceSuite) TestCreateTranc_SuccesLogInfo() {
	log := &loggerStub{}
	service := NewPaymentService(log)
	orderUuid := uuid.New()
	userUuid := uuid.New()
	paymentMethod := model.PaymentMethodCard

	_, err := service.CreateTransaction(s.ctx, orderUuid, userUuid, paymentMethod)
	s.NoError(err)
	s.Equal("Payment completed", log.message)
}
