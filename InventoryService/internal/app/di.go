package app

import (
	api "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/api/inventory/v1"
	repository "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/part"
	service "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service/part"
	inventoryV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
)

type container struct {
	database       *mongo.Database
	collectionName string
	repository     *repository.MongoRepo
	service        *service.InventoryService
	api            inventoryV1.InventoryServiceServer
}

func newContainer(database *mongo.Database, collectionName string) *container {
	return &container{database: database, collectionName: collectionName}
}

func (c *container) Repository() *repository.MongoRepo {
	if c.repository == nil {
		c.repository = repository.NewMongoRepo(c.database, c.collectionName)
	}
	return c.repository
}

func (c *container) Service() *service.InventoryService {
	if c.service == nil {
		c.service = service.NewInventoryService(c.Repository())
	}
	return c.service
}

func (c *container) API() inventoryV1.InventoryServiceServer {
	if c.api == nil {
		c.api = api.NewInventoryServer(c.Service())
	}
	return c.api
}
