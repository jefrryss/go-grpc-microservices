package service

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
)

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
	var buf bytes.Buffer

	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	orderUuid := uuid.New()
	userUuid := uuid.New()
	paymentMethod := model.PaymentMethodCard

	trancUUID, err := s.service.CreateTransaction(s.ctx, orderUuid, userUuid, paymentMethod)
	s.NoError(err)

	text := buf.String()
	expectedLogMessage := fmt.Sprintf("Оплата прошла успешно, transaction_uuid: %s", trancUUID.String())
	s.Contains(text, expectedLogMessage)
}
