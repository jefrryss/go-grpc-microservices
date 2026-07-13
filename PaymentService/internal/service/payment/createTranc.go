package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
	"go.uber.org/zap"
)

func (p *PaymentService) CreateTransaction(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}
	if orderUUID == uuid.Nil {
		return uuid.Nil, model.ErrEmptyOrderUUID
	}
	if userUUID == uuid.Nil {
		return uuid.Nil, model.ErrEmptyUserUUID
	}
	trancUUID := uuid.New()

	p.logger.Info(ctx, "Payment completed", zap.String("transaction_uuid", trancUUID.String()))
	return trancUUID, nil

}
