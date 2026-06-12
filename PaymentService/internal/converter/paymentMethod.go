package converter

import (
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/PaymentService/internal/model"
	payment_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/payment/v1"
)

func ToDomainPaymentMethod(protoMethod payment_v1.PaymentMethod) (model.PaymentMethod, error) {
	switch protoMethod {
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return "", model.ErrInvalidPaymentMethod
	default:
		return "", fmt.Errorf("%w: unknown code %d", model.ErrInvalidPaymentMethod, protoMethod)
	}
}
