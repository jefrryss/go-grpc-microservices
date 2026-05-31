package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/delivery"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/service"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const GrpcPort = 50051
const HttpPort = 8080

func main() {
	log.Println("Starting OrderService initialization...")

	inventoryConn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := delivery.NewGrpcInventoryClient(inventoryConn)
	log.Println("Connected to InventoryService on localhost:50052")

	paymentConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := delivery.NewGrpcPaymentClient(paymentConn)
	log.Println("Connected to PaymentService on localhost:50053")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on TCP port %d: %v", GrpcPort, err)
	}
	log.Printf("TCP listener created on port %d", GrpcPort)

	grpcServer := grpc.NewServer()
	memoryRepo := repository.NewOrderMemory()
	orderService := service.NewOrderService(memoryRepo, paymentClient, inventoryClient)
	orderServer := delivery.NewOrderServer(orderService)

	order_v1.RegisterOrderServiceServer(grpcServer, orderServer)
	reflection.Register(grpcServer)
	log.Println("gRPC server registered successfully")

	go func() {
		log.Printf("gRPC server is running and listening on port :%d", GrpcPort)
		err := grpcServer.Serve(listener)
		if err != nil {
			log.Fatalf("gRPC server failed to serve: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := fmt.Sprintf("localhost:%d", GrpcPort)

	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/redoc.html")
	})
	rootMux.HandleFunc("/swagger/order.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger/order.swagger.json")
	})

	rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
			return
		}
		mux.ServeHTTP(w, r)
	})

	err = order_v1.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}
	serverHttp := &http.Server{
		Addr:    fmt.Sprintf(":%d", HttpPort),
		Handler: rootMux,
	}
	go func() {
		log.Printf("HTTP server (gRPC-Gateway) is running on port :%d", HttpPort)
		if err := serverHttp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to serve: %v", err)
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Shutdown signal received")

	log.Println("Shutting down HTTP server...")
	shutDownCtx, shutDownCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer shutDownCancel()
	if err := serverHttp.Shutdown(shutDownCtx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped successfully")
}
