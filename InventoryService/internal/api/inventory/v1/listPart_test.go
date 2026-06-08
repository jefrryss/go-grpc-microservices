package api

import (
	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *InventoryServerSuite) TestListPart_Success() {
	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Categories: []inventory_v1.Category{
				inventory_v1.Category_CATEGORY_ENGINE,
				inventory_v1.Category_CATEGORY_FUEL,
			},
		},
	}
	partID := uuid.New()
	domainParts := []*model.Part{{PartID: partID, Category: model.CategoryEngine}}
	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return(domainParts, nil).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.NoError(err)

	s.NotNil(res)
	s.Equal(partID.String(), res.Parts[0].Uuid)
}
func (s *InventoryServerSuite) TestListPart_NilFilter() {
	req := &inventory_v1.ListPartsRequest{
		Filter: nil,
	}
	partID := uuid.New()
	domainParts := []*model.Part{{PartID: partID, Category: model.CategoryEngine}}

	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return(domainParts, nil).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.NoError(err)
	s.NotNil(res)

	s.Equal(partID.String(), res.Parts[0].Uuid)
}

func (s *InventoryServerSuite) TestListPart_InvalidFieled() {
	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: []string{"invalid uuid"},
		},
	}
	res, err := s.server.ListParts(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())

	s.serviceMock.AssertNotCalled(s.T(), "ListParts", mock.Anything, mock.Anything)
}

func (s *InventoryServerSuite) TestListPart_ErrServiceNotFound() {
	req := &inventory_v1.ListPartsRequest{}
	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.NotFound, st.Code())
}
func (s *InventoryServerSuite) TestListPart_ErrInvalidCategory() {
	req := &inventory_v1.ListPartsRequest{}
	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return(nil, model.ErrInvalidCategory).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

func (s *InventoryServerSuite) TestListPart_ErrInvalidFilter() {
	req := &inventory_v1.ListPartsRequest{}
	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return(nil, model.ErrInvalidFilter).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

func (s *InventoryServerSuite) TestListPart_EmptyResult() {
	req := &inventory_v1.ListPartsRequest{}
	s.serviceMock.On("ListParts", mock.Anything, mock.Anything).Return([]*model.Part{}, nil).Once()

	res, err := s.server.ListParts(s.ctx, req)
	s.NotNil(res)
	s.NoError(err)

	s.Empty(res.Parts)
}
