package repository

import (
	"context"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	converter "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/converter"
)

func (m *MemoryRepo) GetAll(ctx context.Context) ([]*model.Part, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	res := make([]*model.Part, 0, len(m.data))

	for _, val := range m.data {
		modPart, err := converter.ConvertRepoPartToModelPart(val)
		if err != nil {
			return nil, fmt.Errorf("failed to convert part (ID: %s): %w", val.PartID, err)
		}
		res = append(res, modPart)
	}
	return res, nil
}
