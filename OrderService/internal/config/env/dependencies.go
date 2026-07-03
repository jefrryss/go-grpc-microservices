package env

type DependencyConfig struct {
	inventory string
	payment   string
}

func NewDependencyConfig() *DependencyConfig {
	return &DependencyConfig{
		inventory: valueOrDefault("INVENTORY_SERVICE_ADDRESS", "localhost:50052"),
		payment:   valueOrDefault("PAYMENT_SERVICE_ADDRESS", "localhost:50053"),
	}
}

func (c *DependencyConfig) InventoryAddress() string { return c.inventory }
func (c *DependencyConfig) PaymentAddress() string   { return c.payment }
