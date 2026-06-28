package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	converter "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/converter"
	repoModel "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	var part repoModel.PartRepo
	err := m.collection.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("find part %s: %w", id, err)
	}
	modelPart, err := converter.ConvertRepoPartToModelPart(&part)
	if err != nil {
		return nil, fmt.Errorf("invalid convert part to model: %w", err)
	}
	return modelPart, nil
}
