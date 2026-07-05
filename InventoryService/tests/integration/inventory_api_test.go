//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/app"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/config"
	configMocks "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/config/mocks"
	inventoryV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestInventoryAPIWithMongoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:7.0.5",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForListeningPort("27017/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mongoContainer.Terminate(context.Background())) })

	host, err := mongoContainer.Host(ctx)
	require.NoError(t, err)
	port, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	require.NoError(t, err)
	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, mongoClient.Disconnect(context.Background())) })

	partID := uuid.New().String()
	now := time.Now().UTC()
	_, err = mongoClient.Database("inventory_test").Collection("parts").InsertOne(ctx, bson.M{
		"_id":                  partID,
		"name":                 "Ion Engine",
		"description":          "Integration test engine",
		"price":                1500.0,
		"stock_quantity":       3,
		"category":             "ENGINE",
		"manufacturer_name":    "Test Labs",
		"manufacturer_country": "RU",
		"manufacturer_website": "https://example.test",
		"length":               2.0,
		"width":                1.0,
		"height":               1.0,
		"weight":               100.0,
		"tags":                 []string{"engine"},
		"created_at":           now,
		"updated_at":           now,
	})
	require.NoError(t, err)

	cfg := &config.Config{
		GRPC: configMocks.GRPCConfig{Value: "127.0.0.1:0"},
		Mongo: configMocks.MongoConfig{
			URIValue:        mongoURI,
			DatabaseValue:   "inventory_test",
			CollectionValue: "parts",
		},
		Logger: configMocks.LoggerConfig{LevelValue: "error"},
	}
	application, err := app.New(ctx, cfg)
	require.NoError(t, err)

	appCtx, stopApp := context.WithCancel(ctx)
	go func() { _ = application.Run(appCtx) }()
	t.Cleanup(func() {
		stopApp()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		require.NoError(t, application.Close(shutdownCtx))
	})

	connection, err := grpc.NewClient(application.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })

	client := inventoryV1.NewInventoryServiceClient(connection)
	response, err := client.GetPart(ctx, &inventoryV1.GetPartRequest{Uuid: partID}, grpc.WaitForReady(true))
	require.NoError(t, err)
	require.Equal(t, partID, response.GetPart().GetUuid())
	require.Equal(t, "Ion Engine", response.GetPart().GetName())
}
