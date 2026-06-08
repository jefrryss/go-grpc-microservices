package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *InventiryServiceSuite) TestListParts_Success() {
	part1 := &model.Part{PartID: uuid.New(), Category: model.CategoryEngine, Tags: []string{"metal"}}
	part2 := &model.Part{PartID: uuid.New(), Category: model.CategoryFuel, Tags: []string{"rubber"}}
	allParts := []*model.Part{part1, part2}

	s.repositoryMock.On("GetAll", s.ctx).Return(allParts, nil).Once()

	filter := &model.Filter{Categories: []model.Category{model.CategoryEngine}}
	res, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Len(res, 1)
	s.Equal(model.CategoryEngine, res[0].Category)
}

func (s *InventiryServiceSuite) TestListParts_EmptyResult() {
	allParts := []*model.Part{{PartID: uuid.New(), Category: model.CategoryEngine}}

	s.repositoryMock.On("GetAll", s.ctx).Return(allParts, nil).Once()

	filter := &model.Filter{Categories: []model.Category{model.CategoryFuel}}
	res, err := s.service.ListParts(s.ctx, filter)

	s.NoError(err)
	s.Empty(res)
}

func (s *InventiryServiceSuite) TestListParts_RepositoryError() {
	s.repositoryMock.On("GetAll", s.ctx).Return(nil, errors.New("db fail")).Once()

	res, err := s.service.ListParts(s.ctx, &model.Filter{})

	s.Error(err)
	s.Nil(res)
	s.Contains(err.Error(), "db fail")
}

func (s *InventiryServiceSuite) TestListParts_CanceledContext() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	res, err := s.service.ListParts(ctx, &model.Filter{})

	s.ErrorIs(err, context.Canceled)
	s.Nil(res)
	s.repositoryMock.AssertNotCalled(s.T(), "GetAll", mock.Anything)
}
