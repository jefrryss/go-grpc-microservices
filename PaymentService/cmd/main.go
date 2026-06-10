package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/delivery"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort int = 50053

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		fmt.Printf("Error started tcp port: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	paymentService := service.NewPaymentService()
	paymentServer := delivery.NewPaymentServer(paymentService)

	payment_v1.RegisterPaymentServiceServer(grpcServer, paymentServer)

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
