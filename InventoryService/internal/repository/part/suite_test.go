package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func toBSONDocument(t *testing.T, value any) bson.D {
	t.Helper()

	data, err := bson.Marshal(value)
	require.NoError(t, err)

	var document bson.D
	require.NoError(t, bson.Unmarshal(data, &document))

	return document
}
