package repository

import (
	inventory_v1 "InventoryService/pkg/inventory/v1"
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNotFound = errors.New("repository: item not found")
)

type Repo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*inventory_v1.Part, error)
	GetAll(ctx context.Context) []*inventory_v1.Part
}

type MemoryRepo struct {
	data map[uuid.UUID]*inventory_v1.Part
	rw   sync.RWMutex
}

func NewMemoryRepo() Repo {
	m := &MemoryRepo{
		data: make(map[uuid.UUID]*inventory_v1.Part),
	}

	testID := uuid.New()
	m.data[testID] = &inventory_v1.Part{
		Uuid:  testID.String(),
		Name:  "Тестовый двигатель",
		Price: 100500,
	}

	return m
}

func (m *MemoryRepo) GetByID(ctx context.Context, id uuid.UUID) (*inventory_v1.Part, error) {
	m.rw.RLock()
	defer m.rw.RUnlock()

	if val, ok := m.data[id]; ok {
		clone := proto.Clone(val).(*inventory_v1.Part)
		return clone, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryRepo) GetAll(ctx context.Context) []*inventory_v1.Part {
	m.rw.RLock()
	defer m.rw.RUnlock()

	res := make([]*inventory_v1.Part, 0, len(m.data))

	for _, val := range m.data {
		clone := proto.Clone(val).(*inventory_v1.Part)
		res = append(res, clone)
	}
	return res
}
