package app

import (
	api "github.com/jefrryss/go-grpc-microservices/PaymentService/internal/api/payment/v1"
	service "github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service/payment"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	paymentV1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
)

type container struct {
	service *service.PaymentService
	api     paymentV1.PaymentServiceServer
	logger  *logger.Logger
}

func newContainer(log *logger.Logger) *container {
	return &container{logger: log}
}

func (c *container) Service() *service.PaymentService {
	if c.service == nil {
		c.service = service.NewPaymentService(c.logger)
	}
	return c.service
}

func (c *container) API() paymentV1.PaymentServiceServer {
	if c.api == nil {
		c.api = api.NewPaymentServer(c.Service())
	}
	return c.api
}
