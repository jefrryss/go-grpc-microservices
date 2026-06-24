package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoRepoGetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		expected := testPart(uuid.New(), "Hyperion Engine V8", model.CategoryEngine)
		namespace := mt.DB.Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			namespace,
			mtest.FirstBatch,
			toBSONDocument(mt.T, expected),
		))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		part, err := repo.GetByID(context.Background(), expected.PartID)

		require.NoError(mt, err)
		require.NotNil(mt, part)
		require.Equal(mt, expected.PartID, part.PartID)
		require.Equal(mt, expected.Name, part.Name)
	})

	mt.Run("not found", func(mt *mtest.T) {
		namespace := mt.DB.Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(0, namespace, mtest.FirstBatch))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		part, err := repo.GetByID(context.Background(), uuid.New())

		require.Nil(mt, part)
		require.ErrorIs(mt, err, model.ErrNotFound)
	})

	mt.Run("conversion error", func(mt *mtest.T) {
		expected := testPart(uuid.New(), "Broken part", model.CategoryEngine)
		expected.Metadata = []byte("invalid json")
		namespace := mt.DB.Name() + "." + mt.Coll.Name()
		mt.AddMockResponses(mtest.CreateCursorResponse(
			0,
			namespace,
			mtest.FirstBatch,
			toBSONDocument(mt.T, expected),
		))

		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		part, err := repo.GetByID(context.Background(), expected.PartID)

		require.Nil(mt, part)
		require.ErrorContains(mt, err, "failed to unmarshal metadata")
	})

	mt.Run("database error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123,
			Message: "database unavailable",
		}))

		id := uuid.New()
		repo := NewMongoRepo(mt.DB, mt.Coll.Name())
		part, err := repo.GetByID(context.Background(), id)

		require.Nil(mt, part)
		require.ErrorContains(mt, err, "find part")
	})
}
