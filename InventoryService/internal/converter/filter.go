package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func ConvertToDomainFilter(filter *inventory_v1.PartsFilter) (*model.Filter, error) {
	if filter == nil {
		return nil, nil
	}

	var arrId = make([]uuid.UUID, 0, len(filter.Uuids))
	for _, idStr := range filter.Uuids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid filter UUID '%s': %w", idStr, err)
		}
		arrId = append(arrId, id)
	}

	var categories []model.Category
	for _, protoCat := range filter.Categories {
		domCat, err := ConvertProtoCatToModelCat(protoCat)
		if err != nil {
			return nil, fmt.Errorf("invalid filter category: %w", err)
		}
		categories = append(categories, domCat)

	}
	return &model.Filter{
		Uuids:      arrId,
		Names:      filter.Names,
		Categories: categories,
		Countries:  filter.ManufacturerCountries,
		Tags:       filter.Tags,
	}, nil
}
