package spec

import "time"

type ProbeSpec struct {
	Name string
	Target
	Probe
}

// Only one of the field should be set.
type Target struct {
	Namespace string

	Pod *string
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

	StatusCode  int
	StatusRegex *string
	BodyRegex   *string
	LineMatch   *string
}

type PingProbe struct {
	Interval time.Duration
}
