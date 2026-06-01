package converter

import (
	"fmt"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func ConvertProtoCatToModelCat(category inventory_v1.Category) (model.Category, error) {
	switch category {
	case inventory_v1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine, nil
	case inventory_v1.Category_CATEGORY_FUEL:
		return model.CategoryFuel, nil
	case inventory_v1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole, nil
	case inventory_v1.Category_CATEGORY_UNSPECIFIED:
		return "", nil
	default:
		return "", fmt.Errorf("%w: unknown category code %d", model.ErrInvalidCategory, category)
	}
}
func ToProtoCategory(domainCategory model.Category) inventory_v1.Category {
	switch domainCategory {
	case model.CategoryEngine:
		return inventory_v1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventory_v1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventory_v1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventory_v1.Category_CATEGORY_WING
	default:
		return inventory_v1.Category_CATEGORY_UNSPECIFIED
	}
}
