package reporting

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusReporter struct {
	stateGauge prometheus.Gauge
}

func (pr *PrometheusReporter) ReportState(s int) {
	pr.stateGauge.Set(float64(s))
}

func NewPrometheus(proberName string) *PrometheusReporter {
	sg := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "kprober",
			Subsystem: proberName,
			Name:      "state",
			Help:      "The state of the prober.",
		})

	prometheus.MustRegister(sg)

	return &PrometheusReporter{
		stateGauge: sg,
	}
}
