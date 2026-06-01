package repository

import (
	"encoding/json"
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	repository "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/repository/model"
)

func ConvertModelPartToRepoPart(part *model.Part) (*repository.PartRepo, error) {
	if part == nil {
		return nil, model.ErrNilModelPart
	}

	var metaDataBytes []byte
	if part.MetaData != nil {
		bytes, err := json.Marshal(part.MetaData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata in convertor model to repo: %w", err)
		}
		metaDataBytes = bytes
	}

	partRepo := &repository.PartRepo{
		PartID:              part.PartID,
		Name:                part.Name,
		Description:         part.Description,
		Price:               part.Price,
		StockQuantity:       part.StockQuantity,
		Category:            string(part.Category),
		ManufacturerName:    part.Manufacturer.Name,
		ManufacturerCountry: part.Manufacturer.Country,
		ManufacturerWebsite: part.Manufacturer.Website,
		Length:              part.Dimensions.Length,
		Width:               part.Dimensions.Width,
		Height:              part.Dimensions.Height,
		Weight:              part.Dimensions.Weight,
		Tags:                part.Tags,
		Metadata:            metaDataBytes,
		CreatedAt:           part.CreatedAt,
		UpdatedAt:           part.UpdatedAt,
	}

	return partRepo, nil
}

func ConvertRepoPartToModelPart(part *repository.PartRepo) (*model.Part, error) {
	if part == nil {
		return nil, model.ErrNilRepoPArt
	}

	var metaData map[string]any
	if len(part.Metadata) > 0 {
		if err := json.Unmarshal(part.Metadata, &metaData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	modelPart := &model.Part{
		PartID:        part.PartID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions: model.Dimensions{
			Width:  part.Width,
			Height: part.Height,
			Length: part.Length,
			Weight: part.Weight,
		},
		Manufacturer: model.Manufacturer{
			Name:    part.ManufacturerName,
			Country: part.ManufacturerCountry,
			Website: part.ManufacturerWebsite,
		},
		Tags:      part.Tags,
		MetaData:  metaData,
		CreatedAt: part.CreatedAt,
		UpdatedAt: part.UpdatedAt,
	}
	return modelPart, nil
}
