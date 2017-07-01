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
	// Service IP, Pod IP or IP outside Kubernetes network
	IP string

	Pod        *PodTarget
	Pods       *PodsTarget
	Service    *ServiceTarget
	Deployment *DeploymentTarget
}

type PodTarget struct {
	Namespace string
	Name      string
}

type PodsTarget struct {
	Namespace string
	Selectors map[string]string
}

type DeploymentTarget struct {
	Namespace string
	Name      string
}

type ServiceTarget struct {
	Namespace string
	Name      string
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

	Interval time.Duration

	StatusCode int

	BodyMatchesRegexp      []string
	BodyDoesNotMatchRegexp []string
}

type PingProbe struct {
	Interval time.Duration

	Timeout time.Duration
}
