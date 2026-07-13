package service

import (
	"context"

	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/service"
	"github.com/jefrryss/go-grpc-microservices/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ service.Service = (*PaymentService)(nil)

type Logger interface {
	Info(context.Context, string, ...zap.Field)
}

type PaymentService struct{ logger Logger }

func NewPaymentService(loggers ...Logger) *PaymentService {
	var log Logger = logger.NewNop()
	if len(loggers) > 0 && loggers[0] != nil {
		log = loggers[0]
	}
	return &PaymentService{logger: log}
}
