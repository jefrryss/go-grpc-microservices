package api

import (
	"context"
	"testing"

	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type PaymentServerSuite struct {
	suite.Suite
	ctx         context.Context
	serviceMock *mocks.Service
	server      *PaymentServer
}

func (s *PaymentServerSuite) SetupTest() {
	s.ctx = context.Background()
	mockService := mocks.NewService(s.T())
	s.server = NewPaymentServer(mockService)
	s.serviceMock = mockService
}

func (p *PaymentServerSuite) TearDownTest() {
	p.serviceMock.AssertExpectations(p.T())
}

func TestServer(t *testing.T) {
	suite.Run(t, new(PaymentServerSuite))
}
