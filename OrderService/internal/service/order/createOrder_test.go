package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	userUuid := uuid.New()
	idx1 := uuid.New()
	idx2 := uuid.New()
	firstItem := &model.Part{
		ID:    idx1,
		Price: 109,
	}
	secondItem := &model.Part{
		ID:    idx2,
		Price: 123,
	}
	uuids := []uuid.UUID{idx1, idx2}
	s.clInventoryMock.On("GetListParts", mock.Anything, mock.Anything).Return([]*model.Part{secondItem, firstItem}, nil).Once()
	s.repoMock.On("SetOrder", mock.Anything, mock.Anything).Return(nil).Once()

	res, total, err := s.service.CreateOrder(s.ctx, userUuid, uuids)

	s.NoError(err)
	s.Equal(total, firstItem.Price+secondItem.Price)
	s.NotEqual(res, uuid.Nil)

}

func (s *ServiceSuite) TestCreateOrder_CancelCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, total, err := s.service.CreateOrder(ctx, uuid.New(), []uuid.UUID{})
	s.Equal(res, uuid.Nil)
	s.Error(err)
	s.ErrorIs(err, context.Canceled)
	s.Equal(total, 0.0)
	s.clInventoryMock.AssertNotCalled(s.T(), "GetListParts", mock.Anything, mock.Anything)
	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCreateOrder_InvalidClientInventory() {
	errClient := errors.New("client error")
	userID := uuid.New()
	requestedParts := []uuid.UUID{}

	s.clInventoryMock.On("GetListParts", s.ctx, requestedParts).Return(nil, errClient).Once()

	orderID, total, err := s.service.CreateOrder(s.ctx, userID, requestedParts)

	s.Equal(uuid.Nil, orderID)
	s.Equal(float64(0), total)
	s.Error(err)
	s.Equal(errClient.Error(), err.Error())

	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCreateOrder_InvalidAvalibleUUIDS() {
	avalibleUUids := []*model.Part{}
	part := uuid.New()
	requiredUUids := []uuid.UUID{part}

	s.clInventoryMock.On("GetListParts", mock.Anything, mock.Anything).Return(avalibleUUids, nil).Once()
	res, total, err := s.service.CreateOrder(s.ctx, uuid.New(), requiredUUids)
	s.Equal(res, uuid.Nil)
	s.Error(err)
	s.Equal(total, float64(0))
	s.Contains(err.Error(), fmt.Sprintf("part %s not found", part.String()))
	s.repoMock.AssertNotCalled(s.T(), "SetOrder", mock.Anything, mock.Anything)
}

func (s *ServiceSuite) TestCreateOrder_InvalidRepo() {
	partId := uuid.New()
	uuidsArr := []uuid.UUID{partId}
	avaliblePart := []*model.Part{{ID: partId}}

	userID := uuid.New()
	errorRepo := errors.New("error with db")
	s.clInventoryMock.On("GetListParts", mock.Anything, mock.Anything).Return(avaliblePart, nil).Once()
	s.repoMock.On("SetOrder", mock.Anything, mock.Anything).Return(errorRepo).Once()

	res, total, err := s.service.CreateOrder(s.ctx, userID, uuidsArr)
	s.Equal(res, uuid.Nil)
	s.Error(err)
	s.Equal(total, 0.0)
	s.ErrorIs(err, errorRepo)
}
