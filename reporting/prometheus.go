package reporting

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Probe failure
	StateUnknown = State("Unknown")
	// Not reachable
	StateDown = State("Down")
	// Everything is OK
	StateHealthy = State("Healthy")
	// Something is wrong
	StateDegraded = State("Degraded")
)

type State string

type PrometheusReporter struct {
	probeCounter *prometheus.CounterVec
	errorCounter *prometheus.CounterVec
}

func (pr *PrometheusReporter) ReportState(s State, reason string) {
	pr.probeCounter.With(prometheus.Labels{"state": string(s)}).Inc()
	if len(reason) != 0 {
		pr.errorCounter.With(prometheus.Labels{"reason": reason}).Inc()
	}
}

func NewPrometheus(proberName string) *PrometheusReporter {
	pc := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kprober",
			Subsystem: strings.Replace(proberName, "-", "_", -1),
			Name:      "probe_counter",
			Help:      "Total number of probes",
		},
		[]string{"state"},
	)

	errors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kprober",
			Subsystem: strings.Replace(proberName, "-", "_", -1),
			Name:      "error_counter",
			Help:      "Total number of probe errors",
		},
		[]string{"reason"},
	)

	prometheus.MustRegister(pc)
	prometheus.MustRegister(errors)

	return &PrometheusReporter{
		probeCounter: pc,
		errorCounter: errors,
	}
}
