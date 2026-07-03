package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	api "github.com/jefrryss/go-grpc-microservices/OrderService/internal/api/order/v1"
	clientInventory "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/inventory/v1"
	clientPayment "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment/v1"
	repository "github.com/jefrryss/go-grpc-microservices/OrderService/internal/repository/order"
	service "github.com/jefrryss/go-grpc-microservices/OrderService/internal/service/order"
	orderV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	"google.golang.org/grpc"
)

type container struct {
	database      *pgxpool.Pool
	inventoryConn *grpc.ClientConn
	paymentConn   *grpc.ClientConn
	repository    *repository.OrderPostgres
	service       *service.OrderService
	api           orderV1.OrderServiceServer
}

func newContainer(database *pgxpool.Pool, inventoryConn, paymentConn *grpc.ClientConn) *container {
	return &container{database: database, inventoryConn: inventoryConn, paymentConn: paymentConn}
}

func (c *container) Repository() *repository.OrderPostgres {
	if c.repository == nil {
		c.repository = repository.NewOrderPostgres(c.database)
	}
	return c.repository
}

func (c *container) Service() *service.OrderService {
	if c.service == nil {
		c.service = service.NewOrderService(
			c.Repository(),
			clientPayment.NewGrpcPaymentClient(c.paymentConn),
			clientInventory.NewGrpcInventoryClient(c.inventoryConn),
		)
	}
	return c.service
}

func (c *container) API() orderV1.OrderServiceServer {
	if c.api == nil {
		c.api = api.NewOrderServer(c.Service())
	}
	return c.api
}
