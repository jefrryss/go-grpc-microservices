package client

import (
	client "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

var _ client.InventoryClient = (*GrpcInventoryClient)(nil)

type GrpcInventoryClient struct {
	api inventory_v1.InventoryServiceClient
}

func NewGrpcInventoryClient(conn *grpc.ClientConn) *GrpcInventoryClient {
	return &GrpcInventoryClient{
		api: inventory_v1.NewInventoryServiceClient(conn),
	}
}
