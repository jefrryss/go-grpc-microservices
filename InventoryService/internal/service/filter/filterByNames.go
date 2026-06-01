package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

func (i *InventoryFilter) filterByNames(part *model.Part) bool {
	if len(i.allowedNames) == 0 {
		return true
	}
	_, ok := i.allowedNames[part.Name]
	return ok
}
