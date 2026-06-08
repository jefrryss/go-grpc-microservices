package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	repository *MemoryRepo
	ctx        context.Context
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repository = NewMemoryRepo()
}

func (i *RepositorySuite) TearDownTest() {
}

func TestService(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (s *RepositorySuite) seedTestParts() (uuid.UUID, uuid.UUID) {
	id1 := uuid.New()
	id2 := uuid.New()

	part1 := &repoModel.PartRepo{
		PartID:        id1,
		Name:          "Hyperion Engine V8",
		Description:   "Powerful engine for space flights",
		Price:         150000.50,
		StockQuantity: 10,
		Category:      string(model.CategoryEngine),

		Width:  2.5,
		Height: 3.0,
		Length: 4.5,
		Weight: 1200.0,

		ManufacturerName:    "Stark Industries",
		ManufacturerCountry: "USA",
		ManufacturerWebsite: "stark.com",

		Tags:      []string{"engine", "v8", "space"},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	part2 := &repoModel.PartRepo{
		PartID:        id2,
		Name:          "Quantum Porthole",
		Description:   "Unbreakable glass porthole",
		Price:         5000.00,
		StockQuantity: 50,
		Category:      string(model.CategoryPorthole),

		Width:  1.0,
		Height: 1.0,
		Length: 0.1,
		Weight: 45.5,

		ManufacturerName:    "GlassCorp",
		ManufacturerCountry: "Germany",
		ManufacturerWebsite: "glasscorp.de",

		Tags:      []string{"glass", "window"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.repository.data[id1] = part1
	s.repository.data[id2] = part2

	return id1, id2
}
