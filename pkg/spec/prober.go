package spec

import (
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
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProberSpec   `json:"spec"`
	Status            ProberStatus `json:"status"`
}

type ProberSpec struct {
	Target Target
	Probe  Probe
}

type ProberStatus struct {
}

// Only one of the field should be set.
type Target struct {
	Namespace string
	Pod       string
}

// Only one of the field should be set.
type Probe struct {
	HTTP *HTTPProbe
	Ping *PingProbe
}

type HTTPProbe struct {
	Method string // Only Get and Head are supported currently
	Scheme string
	Port   string
	Path   string

	IntervalInMS int

	StatusCode  int
	StatusRegex string
	BodyRegex   string
	LineMatch   string
}

type PingProbe struct {
	IntervalInMS int
}
