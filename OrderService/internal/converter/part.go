package converter

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jefrryss/go-grpc-microservices/OrderService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func ToDomainPart(protoPart *inventory_v1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, model.ErrNilProtoPart
	}

	partID, err := uuid.Parse(protoPart.Uuid)
	if err != nil {
		return nil, fmt.Errorf("invalid part uuid from inventory: %w", err)
	}

	return &model.Part{
		ID:    partID,
		Price: protoPart.Price,
	}, nil
}
