package repository

import (
	"time"

	"github.com/google/uuid"
)

type PartRepo struct {
	PartID              uuid.UUID `bson:"_id"`
	Name                string    `bson:"name"`
	Description         string    `bson:"description"`
	Price               float64   `bson:"price"`
	StockQuantity       int       `bson:"stock_quantity"`
	Category            string    `bson:"category"`
	ManufacturerName    string    `bson:"manufacturer_name"`
	ManufacturerCountry string    `bson:"manufacturer_country"`
	ManufacturerWebsite string    `bson:"manufacturer_website"`
	Length              float64   `bson:"length"`
	Width               float64   `bson:"width"`
	Height              float64   `bson:"height"`
	Weight              float64   `bson:"weight"`
	Tags                []string  `bson:"tags"`
	Metadata            []byte    `bson:"metadata,omitempty"`
	CreatedAt           time.Time `bson:"created_at"`
	UpdatedAt           time.Time `bson:"updated_at"`
}
