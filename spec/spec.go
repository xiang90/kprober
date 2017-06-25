package spec

type ProbeSpec struct {
	Name string
	Target
	Probe
}

// Only one of the field should be set.
type Target struct {
	Pod *string
}

// Only one of the field should be set.
type Probe struct {
	HTTP *HTTPProbe
	Ping *PingProbe
}

type HTTPProbe struct {
	Method string // Only Get and Head are supported currently
	Port   string
	Path   string

	StatusRegex *string
	BodyRegex   *string
	LineMatch   *string
}

type PingProbe struct{}
