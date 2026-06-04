package client

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/converter"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
)

func (g *GrpcPaymentClient) PayOrder(ctx context.Context, orderID, userID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	if err := ctx.Err(); err != nil {
		return uuid.Nil, err
	}

	pbPaymentMethod := converter.ToProtoPaymentMethod(method)

	req := &payment_v1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		UserUuid:      userID.String(),
		PaymentMethod: pbPaymentMethod,
	}

	resp, err := g.api.PayOrder(ctx, req)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUUID, err := uuid.Parse(resp.TransactionUuid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid transaction uuid received from payment service: %w", err)
	}

	return transactionUUID, nil
}
