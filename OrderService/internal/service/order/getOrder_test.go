package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
)

func (s *ServiceSuite) TestGetOrder_Success() {
	orderID := uuid.New()
	expectedOrder := &model.Order{ID: orderID}

	s.repoMock.On("GetOrder", s.ctx, orderID).Return(expectedOrder, nil).Once()

	res, err := s.service.GetOrder(s.ctx, orderID)

	s.NotNil(res)
	s.NoError(err)
	s.Equal(orderID, res.ID)
}

func (s *ServiceSuite) TestGetOrder_CancelCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := s.service.GetOrder(ctx, uuid.New())

	s.Nil(res)
	s.ErrorIs(err, context.Canceled)
	s.repoMock.AssertNotCalled(s.T(), "GetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestGetOrder_RepoError() {
	orderID := uuid.New()
	s.repoMock.On("GetOrder", s.ctx, orderID).Return(nil, model.ErrOrderNotFound).Once()

	res, err := s.service.GetOrder(s.ctx, orderID)

	s.Nil(res)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
