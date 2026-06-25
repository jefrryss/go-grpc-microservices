package service

import (
	"testing"

	"github.com/jefrryss/go-grpc-microservices/InventoryService/internal/model"
	"github.com/stretchr/testify/require"
)

func TestInventoryFilterByCountry(t *testing.T) {
	part := &model.Part{
		Manufacturer: model.Manufacturer{Country: "USA"},
	}

	require.True(t, NewInventoryFilter(nil).FilterPart(part))
	require.True(t, NewInventoryFilter(&model.Filter{
		Countries: []string{"USA"},
	}).FilterPart(part))
	require.False(t, NewInventoryFilter(&model.Filter{
		Countries: []string{"Germany"},
	}).FilterPart(part))
}
