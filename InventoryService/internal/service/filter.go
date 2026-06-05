package service

import "github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"

type PartFilter interface {
	FilterPart(part *model.Part) bool
}
