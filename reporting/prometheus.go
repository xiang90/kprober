package reporting

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Probe failure
	StateUnknown = State(-2)
	// Not reachable
	StateDown = State(-1)
	// Everything is OK
	StateHealthy = State(0)
	// Something is wrong
	StateDegraded = State(1)
)

type State int

type PrometheusReporter struct {
	stateGauge   prometheus.Gauge
	errorCounter *prometheus.CounterVec
}

func (pr *PrometheusReporter) ReportState(s State, reason string) {
	pr.stateGauge.Set(float64(s))
	if len(reason) != 0 {
		pr.errorCounter.With(prometheus.Labels{"reason": reason}).Inc()
	}
}

func NewPrometheus(proberName string) *PrometheusReporter {
	sg := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "kprober",
			Subsystem: strings.Replace(proberName, "-", "_", -1),
			Name:      "state",
			Help:      "The state of the prober.",
		},
	)

	errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kprober",
			Subsystem: strings.Replace(proberName, "-", "_", -1),
			Name:      "error",
			Help:      "The state of the prober.",
		},
		[]string{"reason"},
	)

	prometheus.MustRegister(sg)
	prometheus.MustRegister(errors)

	return &PrometheusReporter{
		stateGauge:   sg,
		errorCounter: errors,
	}
}
