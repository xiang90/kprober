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
	IP      string         `json:"ip,omitempty"`
	Service *ServiceTarget `json:"service,omitempty"`
}

type ServiceTarget struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// Only one of the field should be set.
type Probe struct {
	HTTP      *HTTPProbe      `json:"http,omitempty"`
	Ping      *PingProbe      `json:"ping,omitempty"`
	Container *ContainerProbe `json:"container,omitempty"`
}

type HTTPProbe struct {
	// Only Get and Head are supported currently
	Method string `json:"method"`
	Scheme string `json:"scheme"`
	Port   int    `json:"port"`
	Path   string `json:"path"`

	Interval time.Duration

	TimeoutSeconds int64 `json:"timeoutSeconds,omitempty"`

	StatusCode int `json:"statusCode"`

	BodyMatchesRegexp      []string
	BodyDoesNotMatchRegexp []string
}

type PingProbe struct {
	PeriodSeconds  int64 `json:"periodSeconds,omitempty"`
	TimeoutSeconds int64 `json:"timeoutSeconds,omitempty"`
}

// ContainerProbe specifies a container that can probe.
//
// The container MUST can execute command `probe`.
// It MUST write the probe result to stdout in the format `state result\n`.
// state MUST be an integer and reason must be a human-readable string.
//
// The container SHOULD accept environment variable IP, which contains the
// IP address (which is generated from target spec by the operator) it should probe.
// The container MAY accept environment variable Target, which contains the JSON
// format of the target spec it should probe.
type ContainerProbe struct {
	Image string `json:"image,omitempty"`
}
