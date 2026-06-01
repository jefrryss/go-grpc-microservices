package model

import (
	"time"

	"github.com/google/uuid"
)

type Category string

const (
	CategoryEngine   Category = "Engine"
	CategoryFuel     Category = "Fuel"
	CategoryPorthole Category = "Porthole"
	CategoryWing     Category = "Wing"
)

type Part struct {
	PartID        uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	MetaData      map[string]any
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Dimensions struct {
	Width  float64
	Height float64
	Length float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}
