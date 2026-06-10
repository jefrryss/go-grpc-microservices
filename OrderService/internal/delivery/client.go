package delivery

import (
	"context"

	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/domain"

	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type GrpcInventoryClient struct {
	api inventory_v1.InventoryServiceClient
}

func NewGrpcInventoryClient(conn *grpc.ClientConn) *GrpcInventoryClient {
	return &GrpcInventoryClient{
		api: inventory_v1.NewInventoryServiceClient(conn),
	}
}

type GrpcPaymentClient struct {
	api payment_v1.PaymentServiceClient
}

func NewGrpcPaymentClient(conn *grpc.ClientConn) *GrpcPaymentClient {
	return &GrpcPaymentClient{
		api: payment_v1.NewPaymentServiceClient(conn),
	}
}

func (g *GrpcPaymentClient) PayOrder(ctx context.Context, orderUUID string, userUUID string, paymentMethod domain.PaymentMethod) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	var pbPaymentMethod payment_v1.PaymentMethod
	switch paymentMethod {
	case domain.PaymentMethodCard:
		pbPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case domain.PaymentMethodSBP:
		pbPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case domain.PaymentMethodCreditCard:
		pbPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case domain.PaymentMethodInvestorMoney:
		pbPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		pbPaymentMethod = payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	req := &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: pbPaymentMethod,
	}

	resp, err := g.api.PayOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}

func (g *GrpcInventoryClient) GetListParts(ctx context.Context, parts []string) ([]*inventory_v1.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	req := &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: parts,
		},
	}
	res, err := g.api.ListParts(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Parts, nil
}

func (g *GrpcInventoryClient) GetPart(ctx context.Context, uuidPart string) (*inventory_v1.Part, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req := &inventory_v1.GetPartRequest{
		Uuid: uuidPart,
	}
	res, err := g.api.GetPart(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.Part, nil
}
