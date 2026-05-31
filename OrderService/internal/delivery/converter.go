package delivery

import (
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/domain"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
)

func ToDomainPaymentMethod(pbMethod order_v1.PaymentMethod) domain.PaymentMethod {
	switch pbMethod {
	case order_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return domain.PaymentMethodCard
	case order_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return domain.PaymentMethodSBP
	case order_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return domain.PaymentMethodCreditCard
	case order_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return domain.PaymentMethodInvestorMoney
	default:
		return domain.PaymentMethodUnknown
	}
}

func ToProtoPaymentMethod(domainMethod domain.PaymentMethod) order_v1.PaymentMethod {
	switch domainMethod {
	case domain.PaymentMethodCard:
		return order_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case domain.PaymentMethodSBP:
		return order_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case domain.PaymentMethodCreditCard:
		return order_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case domain.PaymentMethodInvestorMoney:
		return order_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return order_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

func ToProtoOrder(order *domain.Order) *order_v1.Order {
	if order == nil {
		return nil
	}

	partUUIDS := make([]string, 0, len(order.PartIDs))
	for _, id := range order.PartIDs {
		partUUIDS = append(partUUIDS, id.String())
	}

	var transactionID string
	if order.TransactionID.Valid {
		transactionID = order.TransactionID.UUID.String()
	}

	return &order_v1.Order{
		OrderUuid:       order.ID.String(),
		UserUuid:        order.UserID.String(),
		PartUuids:       partUUIDS,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: transactionID,
		PaymentMethod:   ToProtoPaymentMethod(order.PaymentMethod),
		Status:          string(order.Status),
	}
}
