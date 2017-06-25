# kprober
prober as a service

```go

type ProberSpec struct {
    Name string
    Target
    Probe
}

type ProberStatus struct {
    PrometheusEndpoint string
}

// Only one of the field should be set.
type Target struct {
    Pod *string
    Service *string
    IP *string
    ReplicaSet *string
    Deployment *string
}

// Only one of the field should be set.
type Probe struct {
    HTTP *HTTPProbe
    TCP  *TCPProbe
    Ping *PingProbe
    Exec *ExecProbe
}

type HTTPProbe struct {
    StatusRegex *string
    BodyRegex *string
    LineMatch *string
}

type TCPProbe struct {
    StreamRegex *string
}

type PingProbe struct {}
```

Expose the probing result as a Prometheus metrics

```
prober_{Name} OK
```
