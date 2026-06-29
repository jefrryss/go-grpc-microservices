package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/jefrryss/go-grpc-microservices/OrderService/internal/api/order/v1"
	clientInventory "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory/v1"
	clientPayment "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment/v1"
	repository "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/order"
	service "github.com/jefrryss/go-grpc-microservices/OrderService/internal/service/order"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const GrpcPort = 50051
const HttpPort = 8080

func main() {
	log.Println("Starting OrderService initialization...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	inventoryServiceAddress := os.Getenv("INVENTORY_SERVICE_ADDRESS")
	if inventoryServiceAddress == "" {
		inventoryServiceAddress = "localhost:50052"
	}

	paymentServiceAddress := os.Getenv("PAYMENT_SERVICE_ADDRESS")
	if paymentServiceAddress == "" {
		paymentServiceAddress = "localhost:50053"
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create PostgreSQL pool: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	inventoryConn, err := grpc.NewClient(inventoryServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to InventoryService: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := clientInventory.NewGrpcInventoryClient(inventoryConn)
	log.Printf("InventoryService client configured for %s", inventoryServiceAddress)

	paymentConn, err := grpc.NewClient(paymentServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to PaymentService: %v", err)
	}
	defer paymentConn.Close()
	paymentClient := clientPayment.NewGrpcPaymentClient(paymentConn)
	log.Printf("PaymentService client configured for %s", paymentServiceAddress)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", GrpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on TCP port %d: %v", GrpcPort, err)
	}
	log.Printf("TCP listener created on port %d", GrpcPort)

	grpcServer := grpc.NewServer()
	orderRepository := repository.NewOrderPostgres(db)
	orderService := service.NewOrderService(orderRepository, paymentClient, inventoryClient)
	orderServer := api.NewOrderServer(orderService)

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

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{UseProtoNames: true},
		}),
		runtime.WithForwardResponseOption(forwardHTTPStatus),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := fmt.Sprintf("localhost:%d", GrpcPort)

	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/redoc.html")
	})
	rootMux.HandleFunc("/swagger/order.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swagger/order.swagger.json")
	})
	rootMux.HandleFunc("/openapi/order.openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi/order.openapi.yaml")
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

func forwardHTTPStatus(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	serverMetadata, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	values := serverMetadata.HeaderMD.Get("x-http-code")
	if len(values) == 0 {
		return nil
	}

	statusCode, err := strconv.Atoi(values[0])
	if err != nil {
		return fmt.Errorf("parse HTTP status code: %w", err)
	}

	delete(serverMetadata.HeaderMD, "x-http-code")
	w.Header().Del("Grpc-Metadata-X-Http-Code")
	w.WriteHeader(statusCode)

	return nil
}
