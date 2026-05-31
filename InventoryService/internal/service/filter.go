package service

import (
	"github.com/google/uuid"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

type InventoryFilter struct {
	filters      *inventory_v1.PartsFilter
	allowedUUIDs map[uuid.UUID]struct{}
	allowedNames map[string]struct{}
	allowedTags  map[string]struct{}
}

func NewInventoryFilter(filters *inventory_v1.PartsFilter) *InventoryFilter {
	allowedUUID, allowedNames, allowedTags := createAllowedFeatures(filters)
	return &InventoryFilter{
		filters:      filters,
		allowedUUIDs: allowedUUID,
		allowedNames: allowedNames,
		allowedTags:  allowedTags,
	}
}

func createAllowedFeatures(filters *inventory_v1.PartsFilter) (map[uuid.UUID]struct{}, map[string]struct{}, map[string]struct{}) {

	if filters == nil {
		return nil, nil, nil
	}

	mapUUID := make(map[uuid.UUID]struct{}, len(filters.Uuids))
	mapNames := make(map[string]struct{}, len(filters.Names))
	mapTags := make(map[string]struct{}, len(filters.Tags))

	if len(filters.Uuids) > 0 {
		for _, val := range filters.Uuids {
			uuid, err := uuid.Parse(val)
			if err == nil {
				mapUUID[uuid] = struct{}{}
			}
		}
	}
	if len(filters.Names) > 0 {
		for _, val := range filters.Names {
			mapNames[val] = struct{}{}
		}
	}
	if len(filters.Tags) > 0 {
		for _, val := range filters.Tags {
			mapTags[val] = struct{}{}
		}
	}
	return mapUUID, mapNames, mapTags
}

func (i *InventoryFilter) FilterPart(part *inventory_v1.Part) bool {
	if !i.filterByUUID(part) {
		return false
	}
	if !i.filterByNames(part) {
		return false
	}
	if !i.filterByTags(part) {
		return false
	}
	if !i.filterByCategory(part) {
		return false
	}
	if !i.filterByCountry(part) {
		return false
	}
	return true
}

func (i *InventoryFilter) filterByUUID(part *inventory_v1.Part) bool {
	if len(i.allowedUUIDs) == 0 {
		return true
	}
	uuid, err := uuid.Parse(part.Uuid)
	if err != nil {
		return false
	}

	_, ok := i.allowedUUIDs[uuid]
	return ok
}

func (i *InventoryFilter) filterByNames(part *inventory_v1.Part) bool {
	if len(i.allowedNames) == 0 {
		return true
	}
	_, ok := i.allowedNames[part.Name]
	return ok
}

func (i *InventoryFilter) filterByTags(part *inventory_v1.Part) bool {
	if len(i.allowedTags) == 0 {
		return true
	}
	for _, val := range part.Tags {
		if _, ok := i.allowedTags[val]; ok {
			return true
		}
	}
	return false
}

func (i *InventoryFilter) filterByCategory(part *inventory_v1.Part) bool {
	if i.filters == nil || len(i.filters.Categories) == 0 {
		return true
	}

	for _, requiredCategory := range i.filters.Categories {
		if part.Category == requiredCategory {
			return true
		}
	}
	return false
}

func (i *InventoryFilter) filterByCountry(part *inventory_v1.Part) bool {
	if i.filters != nil || len(i.filters.ManufacturerCountries) == 0 {
		return true
	}
	if part.Manufacturer == nil {
		return false
	}
	for _, country := range i.filters.ManufacturerCountries {
		if country == part.Manufacturer.Country {
			return true
		}
	}
	return false
}
