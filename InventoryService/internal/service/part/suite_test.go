package service

import (
	"context"
	"testing"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type InventiryServiceSuite struct {
	ctx context.Context
	suite.Suite
	repositoryMock *mocks.Repository
	service        *InventoryService
}

func (i *InventiryServiceSuite) SetupTest() {
	i.ctx = context.Background()
	i.repositoryMock = mocks.NewRepository(i.T())
	i.service = NewInventoryService(i.repositoryMock)
}

func (i *InventiryServiceSuite) TearDownTest() {
	i.repositoryMock.AssertExpectations(i.T())
}

func TestService(t *testing.T) {
	suite.Run(t, new(InventiryServiceSuite))
}
