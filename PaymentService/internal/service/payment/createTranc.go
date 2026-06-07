package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
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

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", trancUUID)
	return trancUUID, nil

}
