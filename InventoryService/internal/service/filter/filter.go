package service

import (
	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
)

type InventoryFilter struct {
	filters      *model.Filter
	allowedUUIDs map[uuid.UUID]struct{}
	allowedNames map[string]struct{}
	allowedTags  map[string]struct{}
}

func NewInventoryFilter(filters *model.Filter) *InventoryFilter {
	allowedUUID, allowedNames, allowedTags := createAllowedFeatures(filters)
	return &InventoryFilter{
		filters:      filters,
		allowedUUIDs: allowedUUID,
		allowedNames: allowedNames,
		allowedTags:  allowedTags,
	}
}

func (i *InventoryFilter) FilterPart(part *model.Part) bool {
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
