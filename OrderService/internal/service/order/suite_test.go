package service

import (
	"context"
	"testing"

	mocksInv "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory/mocks"
	mocksPay "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment/mocks"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context
	repoMock        *mocks.Repository
	clPaymentMock   *mocksPay.PaymentClient
	clInventoryMock *mocksInv.InventoryClient
	service         *OrderService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repoMock = mocks.NewRepository(s.T())
	s.clInventoryMock = mocksInv.NewInventoryClient(s.T())
	s.clPaymentMock = mocksPay.NewPaymentClient(s.T())
	s.service = NewOrderService(s.repoMock, s.clPaymentMock, s.clInventoryMock)
}

func (s *ServiceSuite) TearDownTest() {
	s.repoMock.AssertExpectations(s.T())
	s.clPaymentMock.AssertExpectations(s.T())
	s.clInventoryMock.AssertExpectations(s.T())
}
func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
