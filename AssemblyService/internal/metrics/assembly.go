package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Assembly struct{ duration prometheus.Observer }

func NewAssembly(registerer prometheus.Registerer) *Assembly {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "assembly_duration_seconds", Help: "Duration of ship assembly.",
		Buckets: prometheus.DefBuckets,
	})
	registerer.MustRegister(histogram)
	return &Assembly{duration: histogram}
}

func (m *Assembly) AssemblyCompleted(duration time.Duration) {
	m.duration.Observe(duration.Seconds())
}
