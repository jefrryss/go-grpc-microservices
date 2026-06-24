package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoRepoGetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		first := testPart(uuid.New(), "Hyperion Engine V8", model.CategoryEngine)
		second := testPart(uuid.New(), "Quantum Porthole", model.CategoryPorthole)
		namespace := mt.DB.Name() + "." + mt.Coll.Name()

		mt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			namespace,
			mtest.FirstBatch,
			toBSONDocument(mt.T, first),
			toBSONDocument(mt.T, second),
		))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		parts, err := repo.GetAll(context.Background())

		require.NoError(mt, err)
		require.Len(mt, parts, 2)
		require.Equal(mt, first.PartID, parts[0].PartID)
		require.Equal(mt, second.PartID, parts[1].PartID)
	})

	mt.Run("empty", func(mt *mtest.T) {
		namespace := mt.DB.Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, namespace, mtest.FirstBatch))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		parts, err := repo.GetAll(context.Background())

		require.NoError(mt, err)
		require.NotNil(mt, parts)
		require.Empty(mt, parts)
	})

	mt.Run("conversion error", func(mt *mtest.T) {
		part := testPart(uuid.New(), "Broken part", model.CategoryEngine)
		part.Metadata = []byte("invalid json")
		namespace := mt.DB.Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			namespace,
			mtest.FirstBatch,
			toBSONDocument(mt.T, part),
		))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		parts, err := repo.GetAll(context.Background())

		require.Nil(mt, parts)
		require.ErrorContains(mt, err, "failed to unmarshal metadata")
	})

	mt.Run("database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123,
			Message: "database unavailable",
		}))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		parts, err := repo.GetAll(context.Background())

		require.Nil(mt, parts)
		require.ErrorContains(mt, err, "find parts")
	})
}

func testPart(id uuid.UUID, name string, category model.Category) repoModel.PartRepo {
	now := time.Now().UTC().Truncate(time.Millisecond)

	return repoModel.PartRepo{
		PartID:              id,
		Name:                name,
		Description:         "Test part",
		Price:               100,
		StockQuantity:       10,
		Category:            string(category),
		ManufacturerName:    "Test manufacturer",
		ManufacturerCountry: "USA",
		ManufacturerWebsite: "example.com",
		Length:              1,
		Width:               2,
		Height:              3,
		Weight:              4,
		Tags:                []string{"test"},
		Metadata:            []byte(`{"color":"black"}`),
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
