package main

import (
	"InventoryService/internal/delivery"
	"InventoryService/internal/repository"
	"InventoryService/internal/service"
	inventory_v1 "InventoryService/pkg/inventory/v1"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	server := delivery.NewInventoryServer(serv)

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
