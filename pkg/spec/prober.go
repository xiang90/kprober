package spec

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ProberResourceKind   = "Prober"
	ProberResourcePlural = "probers"
)

type ProberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Prober `json:"items"`
}

type Prober struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ProberSpec   `json:"spec"`
	Status            ProberStatus `json:"status"`
}

type ProberSpec struct {
	Target Target `json:"target"`
	Probe  Probe  `json:"probe"`
}

type ProberStatus struct {
}

// Only one of the field should be set.
type Target struct {
	// Service IP, Pod IP or IP outside Kubernetes network
	IP      string `json:"ip,omitempty"`
	Service *ServiceTarget
}

type ServiceTarget struct {
	Namespace string
	Name      string
}

// Only one of the field should be set.
type Probe struct {
	HTTP *HTTPProbe `json:"http,omitempty"`
	Ping *PingProbe `json:"ping,omitempty"`
}

type HTTPProbe struct {
	Method string // Only Get and Head are supported currently
	Scheme string
	Port   string
	Path   string

	Interval time.Duration

	StatusCode int

	BodyMatchesRegexp      []string
	BodyDoesNotMatchRegexp []string
}

type PingProbe struct {
	PeriodSeconds  int64 `json:"periodSeconds,omitempty"`
	TimeoutSeconds int64 `json:"timeoutSeconds,omitempty"`
}
