package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

func (i *InventoryFilter) filterByTags(part *model.Part) bool {
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
