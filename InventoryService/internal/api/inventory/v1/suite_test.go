package api

import (
	"context"
	"testing"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type InventoryServerSuite struct {
	suite.Suite
	ctx         context.Context
	serviceMock *mocks.Service
	server      *InventoryServer
}

func (s *InventoryServerSuite) SetupTest() {
	s.ctx = context.Background()
	s.serviceMock = mocks.NewService(s.T())
	s.server = NewInventoryServer(s.serviceMock)
}

func (i *InventoryServerSuite) TearDownTest() {
	i.serviceMock.AssertExpectations(i.T())
}

func TestService(t *testing.T) {
	suite.Run(t, new(InventoryServerSuite))
}
