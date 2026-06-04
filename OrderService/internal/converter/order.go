package converter

import (
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	order_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/order/v1"
)

func ToProtoOrder(order *model.Order) *order_v1.Order {
	if order == nil {
		return nil
	}

	partUUIDs := make([]string, 0, len(order.PartIDs))
	for _, partID := range order.PartIDs {
		partUUIDs = append(partUUIDs, partID.String())
	}

	var transactionUUID string
	if order.TransactionID.Valid {
		transactionUUID = order.TransactionID.UUID.String()
	}

	return &order_v1.Order{
		OrderUuid:       order.ID.String(),
		UserUuid:        order.UserID.String(),
		PartUuids:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: transactionUUID,

		PaymentMethod: ToProtoOrderPaymentMethod(order.PaymentMethod),

		Status: string(order.Status),
	}
}
