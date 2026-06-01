package repository

import (
	"time"

	"github.com/google/uuid"
)

type PartRepo struct {
	PartID              uuid.UUID
	Name                string
	Description         string
	Price               float64
	StockQuantity       int
	Category            string
	ManufacturerName    string
	ManufacturerCountry string
	ManufacturerWebsite string
	Length              float64
	Width               float64
	Height              float64
	Weight              float64
	Tags                []string
	Metadata            []byte
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
