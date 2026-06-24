package repository

import (
	"context"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	converter "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/converter"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongoRepo) GetAll(ctx context.Context) ([]*model.Part, error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find parts: %w", err)
	}
	defer cursor.Close(ctx)

	var repoParts []repoModel.PartRepo

	if err := cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("decode parts: %w", err)
	}

	modelParts := make([]*model.Part, 0, len(repoParts))
	for i := range repoParts {
		modelPart, err := converter.ConvertRepoPartToModelPart(&repoParts[i])
		if err != nil {
			return nil, fmt.Errorf(
				"convert part %s: %w",
				repoParts[i].PartID,
				err,
			)
		}
		modelParts = append(modelParts, modelPart)
	}
	return modelParts, nil
}
