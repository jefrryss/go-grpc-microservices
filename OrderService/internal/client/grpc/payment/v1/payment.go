package client

import (
	client "github.com/jefrryss/go-grpc-microservices/OrderService/internal/client/grpc/payment"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

var _ client.PaymentClient = (*GrpcPaymentClient)(nil)

type GrpcPaymentClient struct {
	api payment_v1.PaymentServiceClient
}

func NewGrpcPaymentClient(conn *grpc.ClientConn) *GrpcPaymentClient {
	return &GrpcPaymentClient{
		api: payment_v1.NewPaymentServiceClient(conn),
	}
}
