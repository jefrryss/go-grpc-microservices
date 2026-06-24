package repository

import (
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ repository.Repository = (*MongoRepo)(nil)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(db *mongo.Database, collectionName string) *MongoRepo {
	return &MongoRepo{
		collection: db.Collection(collectionName),
	}
}
