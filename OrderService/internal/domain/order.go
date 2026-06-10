package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "UNKNOWN"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	PartIDs       []uuid.UUID
	TotalPrice    float64
	TransactionID uuid.NullUUID
	PaymentMethod PaymentMethod
	Status        OrderStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
