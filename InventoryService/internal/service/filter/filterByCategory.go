package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

func (i *InventoryFilter) filterByCategory(part *model.Part) bool {
	if i.filters == nil || len(i.filters.Categories) == 0 {
		return true
	}

	for _, requiredCategory := range i.filters.Categories {

		if model.Category(part.Category) == requiredCategory {
			return true
		}
	}
	return false
}
