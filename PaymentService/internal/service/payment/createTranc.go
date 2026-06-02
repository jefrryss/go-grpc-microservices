package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
)

func (p *PaymentService) CreateTransaction(ctx context.Context, orderUUID, userrUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}

	trancUUID := uuid.New()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", trancUUID)
	return trancUUID, nil

}
