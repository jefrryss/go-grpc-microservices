package api

import (
	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *InventoryServerSuite) TestGetPart_Success() {
	partId := uuid.New()
	req := &inventory_v1.GetPartRequest{
		Uuid: partId.String(),
	}

	s.serviceMock.On("GetPart", mock.Anything, mock.Anything).Return(&model.Part{PartID: partId}, nil).Once()
	res, err := s.server.GetPart(s.ctx, req)
	s.NoError(err)
	s.NotNil(res)

	s.Equal(partId.String(), res.Part.Uuid)
}

func (s *InventoryServerSuite) TestGetPart_InvalidUuid() {
	req := &inventory_v1.GetPartRequest{
		Uuid: "invalid_uuid",
	}

	res, err := s.server.GetPart(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(st.Code(), codes.InvalidArgument)
	s.serviceMock.AssertNotCalled(s.T(), "GetPart", mock.Anything, mock.Anything)
}

func (s *InventoryServerSuite) TestGetPart_InvalidUuidFromService() {
	req := &inventory_v1.GetPartRequest{Uuid: uuid.New().String()}

	s.serviceMock.On("GetPart", mock.Anything, mock.Anything).Return(nil, nil).Once()
	res, err := s.server.GetPart(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.Internal, st.Code())

}

func (s *InventoryServerSuite) TestGetPart_InvalidNotFound() {
	req := &inventory_v1.GetPartRequest{
		Uuid: uuid.New().String(),
	}

	s.serviceMock.On("GetPart", mock.Anything, mock.Anything).Return(nil, model.ErrNotFound).Once()

	res, err := s.server.GetPart(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.NotFound, st.Code())
}

func (s *InventoryServerSuite) TestGetPart_InvalidCategory() {
	req := &inventory_v1.GetPartRequest{
		Uuid: uuid.New().String(),
	}

	s.serviceMock.On("GetPart", mock.Anything, mock.Anything).Return(nil, model.ErrInvalidCategory).Once()

	res, err := s.server.GetPart(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}

func (s *InventoryServerSuite) TestGetPart_InvalidFilter() {
	req := &inventory_v1.GetPartRequest{
		Uuid: uuid.New().String(),
	}

	s.serviceMock.On("GetPart", mock.Anything, mock.Anything).Return(nil, model.ErrInvalidFilter).Once()

	res, err := s.server.GetPart(s.ctx, req)
	s.Nil(res)
	s.Error(err)

	st, ok := status.FromError(err)
	s.True(ok)
	s.Equal(codes.InvalidArgument, st.Code())
}
