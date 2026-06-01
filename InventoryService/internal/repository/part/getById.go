package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	converter "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/converter"
)

func (m *MemoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	if val, ok := m.data[id]; ok {
		model, err := converter.ConvertRepoPartToModelPart(val)
		if err != nil {
			return nil, err
		}
		return model, nil
	}
	return nil, model.ErrNotFound
}
