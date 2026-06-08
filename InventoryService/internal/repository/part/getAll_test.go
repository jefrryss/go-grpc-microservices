package repository

import (
	"github.com/google/uuid"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
)

func (s *RepositorySuite) TestGetAll_Success() {
	id1, id2 := s.seedTestParts()

	res, err := s.repository.GetAll(s.ctx)

	s.NoError(err)
	s.Len(res, 2)

	foundIDs := make(map[uuid.UUID]bool)
	for _, part := range res {
		foundIDs[part.PartID] = true
	}

	s.True(foundIDs[id1])
	s.True(foundIDs[id2])
}

func (s *RepositorySuite) TestGetAll_ConversionError() {
	badID := uuid.New()

	s.repository.data[badID] = &repoModel.PartRepo{
		PartID:   badID,
		Metadata: []byte(`{ "key": "broken_json... `),
	}

	res, err := s.repository.GetAll(s.ctx)

	s.Nil(res)
	s.Error(err)
	s.Contains(err.Error(), "failed to unmarshal metadata")
}

func (s *RepositorySuite) TestGetAll_Empty() {
	res, err := s.repository.GetAll(s.ctx)

	s.NoError(err)
	s.NotNil(res)
	s.Empty(res)
}
