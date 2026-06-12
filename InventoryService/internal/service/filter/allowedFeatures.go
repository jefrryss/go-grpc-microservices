package service

import (
	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
)

func createAllowedFeatures(filters *model.Filter) (map[uuid.UUID]struct{}, map[string]struct{}, map[string]struct{}) {

	if filters == nil {
		return nil, nil, nil
	}

	mapUUID := make(map[uuid.UUID]struct{}, len(filters.Uuids))
	mapNames := make(map[string]struct{}, len(filters.Names))
	mapTags := make(map[string]struct{}, len(filters.Tags))

	if len(filters.Uuids) > 0 {
		for _, val := range filters.Uuids {
			mapUUID[val] = struct{}{}
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
