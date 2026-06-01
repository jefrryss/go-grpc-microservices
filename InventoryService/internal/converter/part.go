package converter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
)

func ToDomainPart(protoPart *inventory_v1.Part) (*model.Part, error) {
	if protoPart == nil {
		return nil, model.ErrNilProtoPart
	}

	id, err := uuid.Parse(protoPart.Uuid)
	if err != nil {
		return nil, fmt.Errorf("invalid part uuid: %w", err)
	}

	category, err := ConvertProtoCatToModelCat(protoPart.Category)
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]any)
	for k, v := range protoPart.Metadata {
		switch kind := v.Kind.(type) {
		case *inventory_v1.Value_StringValue:
			metadata[k] = kind.StringValue
		case *inventory_v1.Value_Int64Value:
			metadata[k] = kind.Int64Value
		case *inventory_v1.Value_DoubleValue:
			metadata[k] = kind.DoubleValue
		case *inventory_v1.Value_BoolValue:
			metadata[k] = kind.BoolValue
		}
	}

	var dim model.Dimensions
	if protoPart.Dimensions != nil {
		dim = model.Dimensions{
			Length: protoPart.Dimensions.Length,
			Width:  protoPart.Dimensions.Width,
			Height: protoPart.Dimensions.Height,
			Weight: protoPart.Dimensions.Weight,
		}
	}

	var man model.Manufacturer
	if protoPart.Manufacturer != nil {
		man = model.Manufacturer{
			Name:    protoPart.Manufacturer.Name,
			Country: protoPart.Manufacturer.Country,
			Website: protoPart.Manufacturer.Website,
		}
	}

	var createdAt, updatedAt time.Time
	if protoPart.CreatedAt != nil {
		createdAt = protoPart.CreatedAt.AsTime()
	}
	if protoPart.UpdatedAt != nil {
		updatedAt = protoPart.UpdatedAt.AsTime()
	}

	return &model.Part{
		PartID:        id,
		Name:          protoPart.Name,
		Description:   protoPart.Description,
		Price:         protoPart.Price,
		StockQuantity: int(protoPart.StockQuantity),
		Category:      category,
		Dimensions:    dim,
		Manufacturer:  man,
		Tags:          protoPart.Tags,
		MetaData:      metadata,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func ToProtoPart(domainPart *model.Part) (*inventory_v1.Part, error) {
	if domainPart == nil {
		return nil, model.ErrNilDomainPart
	}

	metadata := make(map[string]*inventory_v1.Value)
	for k, v := range domainPart.MetaData {
		switch val := v.(type) {
		case string:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_StringValue{StringValue: val}}
		case int:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_Int64Value{Int64Value: int64(val)}}
		case int32:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_Int64Value{Int64Value: int64(val)}}
		case int64:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_Int64Value{Int64Value: val}}
		case float32:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_DoubleValue{DoubleValue: float64(val)}}
		case float64:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_DoubleValue{DoubleValue: val}}
		case bool:
			metadata[k] = &inventory_v1.Value{Kind: &inventory_v1.Value_BoolValue{BoolValue: val}}
		}
	}

	return &inventory_v1.Part{
		Uuid:          domainPart.PartID.String(),
		Name:          domainPart.Name,
		Description:   domainPart.Description,
		Price:         domainPart.Price,
		StockQuantity: int64(domainPart.StockQuantity),
		Category:      ToProtoCategory(domainPart.Category),
		Dimensions: &inventory_v1.Dimensions{
			Length: domainPart.Dimensions.Length,
			Width:  domainPart.Dimensions.Width,
			Height: domainPart.Dimensions.Height,
			Weight: domainPart.Dimensions.Weight,
		},
		Manufacturer: &inventory_v1.Manufacturer{
			Name:    domainPart.Manufacturer.Name,
			Country: domainPart.Manufacturer.Country,
			Website: domainPart.Manufacturer.Website,
		},
		Tags:      domainPart.Tags,
		Metadata:  metadata,
		CreatedAt: timestamppb.New(domainPart.CreatedAt),
		UpdatedAt: timestamppb.New(domainPart.UpdatedAt),
	}, nil
}
