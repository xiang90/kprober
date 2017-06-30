package reporting

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Not reachable
	StateDown = State(-1)
	// Everything is OK
	StateHealthy = State(0)
	// Something is wrong
	StateDegraded = State(1)
)

type State int

type PrometheusReporter struct {
	stateGauge *prometheus.GaugeVec
}

func (pr *PrometheusReporter) ReportState(s State, reason string) {
	pr.stateGauge.With(prometheus.Labels{"reason": reason}).Set(float64(s))
}

func NewPrometheus(proberName string) *PrometheusReporter {
	sg := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kprober",
			Subsystem: strings.Replace(proberName, "-", "_", -1),
			Name:      "state",
			Help:      "The state of the prober.",
		},
		[]string{"reason"},
	)

	prometheus.MustRegister(sg)

	return &PrometheusReporter{
		stateGauge: sg,
	}
}
