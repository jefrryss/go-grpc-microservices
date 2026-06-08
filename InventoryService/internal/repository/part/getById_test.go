package repository

import (
	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	repository "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
)

func (s *RepositorySuite) TestGetPart_Success() {
	idx1, _ := s.seedTestParts()
	res, err := s.repository.GetByID(s.ctx, idx1)
	s.NoError(err)
	s.NotNil(res)

	s.Equal(res.PartID, idx1)
}

func (s *RepositorySuite) TestGetPart_NotFounPart() {
	uuid := uuid.New()
	res, err := s.repository.GetByID(s.ctx, uuid)
	s.Nil(res)
	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
}

func (s *RepositorySuite) TestGetPart_InvalidMetaData() {
	id := uuid.New()
	s.repository.data[id] = &repository.PartRepo{
		PartID:   id,
		Metadata: []byte("invalid json"),
	}
	res, err := s.repository.GetByID(s.ctx, id)
	s.Nil(res)
	s.Contains(err.Error(), "failed to unmarshal metadata")
}
