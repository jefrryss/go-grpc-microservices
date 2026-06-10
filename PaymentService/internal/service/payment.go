package service

import (
	"context"
	"errors"
	"log"

	payment_v1 "github.com/jefrryss/go-grpc-microservices/PaymentService/pkg/payment/v1"

	"github.com/google/uuid"
)

var (
	ErrUncorrectUUIDOrder = errors.New("Uncorrect UUID order")
	ErrUncorrectUUIDUser  = errors.New("Unccoert UUID user")
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (p *PaymentService) CreateTrunc(ctx context.Context, orderUUID, userrUUID string, paymentMethod payment_v1.PaymentMethod) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	_, err := uuid.Parse(orderUUID)
	if err != nil {
		return "", ErrUncorrectUUIDOrder
	}

	_, err = uuid.Parse(userrUUID)
	if err != nil {
		return "", ErrUncorrectUUIDUser
	}
	trancUUID := uuid.New().String()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", trancUUID)
	return trancUUID, nil

}
