package metrics

import "github.com/prometheus/client_golang/prometheus"

type Order struct {
	total   prometheus.Counter
	revenue prometheus.Counter
}

func NewOrder(registerer prometheus.Registerer) *Order {
	metrics := &Order{
		total: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_total", Help: "Total number of created orders.",
		}),
		revenue: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "orders_revenue_total", Help: "Total revenue from paid orders.",
		}),
	}
	registerer.MustRegister(metrics.total, metrics.revenue)
	return metrics
}

func (m *Order) OrderCreated()            { m.total.Inc() }
func (m *Order) OrderPaid(amount float64) { m.revenue.Add(amount) }
