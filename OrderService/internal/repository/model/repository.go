package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderRepo struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	PartIDs       []uuid.UUID
	TotalPrice    float64
	TransactionID uuid.NullUUID
	PaymentMethod string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
