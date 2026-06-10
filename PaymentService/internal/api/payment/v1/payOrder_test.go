package api

import (
	"errors"

	"github.com/google/uuid"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *PaymentServerSuite) TestPaymentServer_SuccesPayOrder() {
	userUUID := uuid.New()
	orderUUID := uuid.New()
	paymentMethod := payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	expectedTrancUUID := uuid.New()
	req := &payment_v1.PayOrderRequest{
		UserUuid:      userUUID.String(),
		OrderUuid:     orderUUID.String(),
		PaymentMethod: paymentMethod,
	}

	s.serviceMock.On("CreateTransaction", s.ctx, orderUUID, userUUID, mock.Anything).
		Return(expectedTrancUUID, nil)

	response, err := s.server.PayOrder(s.ctx, req)
	s.NoError(err)
	s.NotNil(response)
	s.Equal(response.TransactionUuid, expectedTrancUUID.String())
}

func (s *PaymentServerSuite) TestPaymentServer_InvalidOrderUUIDFormat() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      uuid.New().String(),
		OrderUuid:     "invalid-uuid",
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_EmptyOrderUUID() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      uuid.New().String(),
		OrderUuid:     uuid.Nil.String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_InvalidUserUUIDFormat() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      "invalid-uuid",
		OrderUuid:     uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_EmptyUserUUID() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      uuid.Nil.String(),
		OrderUuid:     uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_UNSPECIFIEDPaymentMethod() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      uuid.New().String(),
		OrderUuid:     uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_UnknownPaymentMethod() {
	req := &payment_v1.PayOrderRequest{
		UserUuid:      uuid.New().String(),
		OrderUuid:     uuid.New().String(),
		PaymentMethod: payment_v1.PaymentMethod(10000),
	}

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	s.serviceMock.AssertNotCalled(s.T(), "CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (s *PaymentServerSuite) TestPaymentServer_InternalServiceError() {
	userUUID := uuid.New()
	orderUUID := uuid.New()
	req := &payment_v1.PayOrderRequest{
		UserUuid:      userUUID.String(),
		OrderUuid:     orderUUID.String(),
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	s.serviceMock.On("CreateTransaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(uuid.Nil, errors.New("database connection failed"))

	res, err := s.server.PayOrder(s.ctx, req)
	s.Error(err)
	s.Nil(res)

	st, _ := status.FromError(err)
	s.Equal(codes.Internal, st.Code())
	s.Contains(st.Message(), "failed to process payment")
}
