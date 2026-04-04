package metrics

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MonotonicGauge struct {
	gauge *prometheus.GaugeVec
	mutex sync.Mutex
	value *float64
}

func NewMonotonicGauge(name string, help string, labels []string) *MonotonicGauge {
	return &MonotonicGauge{
		gauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: fmt.Sprintf("%s_%s", METRICS_NAMESPACE, name),
			Help: help,
		}, labels),
	}
}

func (this *MonotonicGauge) Set(value float64, labels prometheus.Labels) {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.value != nil && value < *this.value {
		return
	}

	this.value = &value
	this.gauge.With(labels).Set(value)
}

func (this *MonotonicGauge) Collector() prometheus.Collector {
	return this.gauge
}
