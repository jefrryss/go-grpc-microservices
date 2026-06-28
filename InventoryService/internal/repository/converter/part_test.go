package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
	"github.com/stretchr/testify/require"
)

func TestPartUUIDRoundTrip(t *testing.T) {
	expectedID := uuid.New()
	part := &model.Part{PartID: expectedID}

	repoPart, err := ConvertModelPartToRepoPart(part)
	require.NoError(t, err)
	require.Equal(t, expectedID.String(), repoPart.PartID)

	domainPart, err := ConvertRepoPartToModelPart(repoPart)
	require.NoError(t, err)
	require.Equal(t, expectedID, domainPart.PartID)
}

func TestConvertRepoPartToModelPartInvalidUUID(t *testing.T) {
	part, err := ConvertRepoPartToModelPart(&repoModel.PartRepo{PartID: "invalid"})

	require.Nil(t, part)
	require.ErrorContains(t, err, "invalid repository part UUID")
}
