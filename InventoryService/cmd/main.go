package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	api "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/api/inventory/v1"
	repository "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/part"
	service "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/service/part"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort int = 50052

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		fmt.Printf("Error started tcp port: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewMemoryRepo()
	serv := service.NewInventoryService(repo)
	server := api.NewInventoryServer(serv)

	inventory_v1.RegisterInventoryServiceServer(grpcServer, server)

	reflection.Register(grpcServer)

	go func() {
		fmt.Printf("Server started on port:%d\n", grpcPort)
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	grpcServer.GracefulStop()
	fmt.Println("\nserver stopped")
}
