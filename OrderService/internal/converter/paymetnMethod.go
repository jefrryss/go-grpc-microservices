package converter

import (
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
)

func ToProtoPaymentMethod(method model.PaymentMethod) payment_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func ToProtoOrderPaymentMethod(method model.PaymentMethod) order_v1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return order_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return order_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return order_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return order_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethodUnknown:
		return order_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	default:
		return order_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN
	}
}

func ToDomainPaymentMethod(protoMethod order_v1.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case order_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case order_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case order_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case order_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	case order_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN:
		return model.PaymentMethodUnknown
	default:
		return model.PaymentMethodUnknown
	}
}
