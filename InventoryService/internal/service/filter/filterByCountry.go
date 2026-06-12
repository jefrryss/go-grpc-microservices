package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

func (i *InventoryFilter) filterByCountry(part *model.Part) bool {
	if i.filters != nil || len(i.filters.Countries) == 0 {
		return true
	}

	for _, country := range i.filters.Countries {
		if country == part.Manufacturer.Country {
			return true
		}
	}
	return false
}
