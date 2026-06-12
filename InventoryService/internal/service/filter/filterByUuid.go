package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

func (i *InventoryFilter) filterByUUID(part *model.Part) bool {
	if len(i.allowedUUIDs) == 0 {
		return true
	}

	_, ok := i.allowedUUIDs[part.PartID]
	return ok
}
