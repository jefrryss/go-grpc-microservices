package converter

import (
	"testing"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	inventory_v1 "github.com/jefrryss/go-grpc-microservices/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/require"
)

func TestConvertProtoCatToModelCat(t *testing.T) {
	tests := []struct {
		name     string
		proto    inventory_v1.Category
		expected model.Category
	}{
		{name: "engine", proto: inventory_v1.Category_CATEGORY_ENGINE, expected: model.CategoryEngine},
		{name: "fuel", proto: inventory_v1.Category_CATEGORY_FUEL, expected: model.CategoryFuel},
		{name: "porthole", proto: inventory_v1.Category_CATEGORY_PORTHOLE, expected: model.CategoryPorthole},
		{name: "wing", proto: inventory_v1.Category_CATEGORY_WING, expected: model.CategoryWing},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			category, err := ConvertProtoCatToModelCat(test.proto)
			require.NoError(t, err)
			require.Equal(t, test.expected, category)
		})
	}
}

func TestConvertProtoCatToModelCatUnknown(t *testing.T) {
	category, err := ConvertProtoCatToModelCat(inventory_v1.Category(999))

	require.Empty(t, category)
	require.ErrorIs(t, err, model.ErrInvalidCategory)
}
