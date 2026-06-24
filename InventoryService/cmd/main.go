package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/api/inventory/v1"
	repository "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/part"
	service "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service/part"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort int = 50052

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		log.Fatal("MONGO_DATABASE is not set")
	}

	collectionName := os.Getenv("MONGO_COLLECTION")
	if collectionName == "" {
		log.Fatal("MONGO_COLLECTION is not set")
	}

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelConnect()

	mongoClient, err := mongo.Connect(connectCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("connect to MongoDB: %v", err)
	}

	if err := mongoClient.Ping(connectCtx, readpref.Primary()); err != nil {
		log.Fatalf("ping MongoDB: %v", err)
	}

	defer func() {
		disconnectCtx, cancelDisconnect := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelDisconnect()

		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Printf("disconnect from MongoDB: %v", err)
		}
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen on port %d: %v", grpcPort, err)
	}
	defer listener.Close()

	grpcServer := grpc.NewServer()
	repo := repository.NewMongoRepo(mongoClient.Database(databaseName), collectionName)
	serv := service.NewInventoryService(repo)
	server := api.NewInventoryServer(serv)

	inventory_v1.RegisterInventoryServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	go func() {
		fmt.Printf("Server started on port:%d\n", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	grpcServer.GracefulStop()
	fmt.Println("server stopped")
}
