package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *InventiryServiceSuite) TestGetPart_CancelContext() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	uuid := uuid.New()

	res, err := s.service.GetPart(ctx, uuid)

	s.Nil(res)
	s.ErrorIs(err, context.Canceled)
	s.repositoryMock.AssertNotCalled(s.T(), "GetByID", mock.Anything, mock.Anything)
}

func (s *InventiryServiceSuite) TestGetPart_Success() {
	uuidPart := uuid.New()
	excpectedPart := &model.Part{PartID: uuidPart}

	s.repositoryMock.On("GetByID", s.ctx, uuidPart).Return(excpectedPart, nil).Once()

	res, err := s.service.GetPart(s.ctx, uuidPart)
	s.NoError(err)
	s.Equal(res, excpectedPart)
}

func (s *InventiryServiceSuite) TestGetPart_NilUuid() {
	uuid := uuid.Nil

	res, err := s.service.GetPart(s.ctx, uuid)
	s.Nil(res)
	s.ErrorIs(err, model.ErrNilIDPart)
	s.repositoryMock.AssertNotCalled(s.T(), "GetByID", mock.Anything, mock.Anything)
}

func (s *InventiryServiceSuite) TestGetPart_ErrNotFoundPartFromRepo() {
	uuid := uuid.New()

	s.repositoryMock.On("GetByID", s.ctx, uuid).Return(nil, model.ErrNotFound).Once()
	res, err := s.service.GetPart(s.ctx, uuid)
	s.Nil(res)
	s.ErrorIs(err, model.ErrNotFound)
}

func (s *InventiryServiceSuite) TestGetPart_NoErrRepoAndNilPart() {
	uuid := uuid.New()
	s.repositoryMock.On("GetByID", s.ctx, uuid).Return(nil, nil).Once()
	res, err := s.service.GetPart(s.ctx, uuid)
	s.Nil(res)
	s.ErrorIs(err, model.ErrNotFound)
}
