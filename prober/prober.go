package prober

import (
	"context"
	"fmt"
	"time"

	"github.com/xiang90/kprober/k8sutil"
	"github.com/xiang90/kprober/probehttp"
	"github.com/xiang90/kprober/reporting"
	"github.com/xiang90/kprober/spec"
)

type Probe interface {
	State() int
}

type Prober struct {
	spec.ProbeSpec
}

func (p *Prober) Start(ctx context.Context) {
	var ip string
	switch {
	case p.Pod != nil:
		k8sutil.IPFromPod(*p.Pod)
	default:
		panic("target unspecified")
	}

	rp := reporting.NewPrometheus(p.Name)

	var probe Probe

	switch {
	case p.HTTP != nil:
		url := fmt.Sprintf("%s://%s:%s/%s", p.HTTP.Scheme, ip, p.HTTP.Port, p.HTTP.Path)
		ph := &probehttp.Probe{
			URL:    url,
			Method: p.HTTP.Method,
		}
		ph.Start(context.TODO())
		probe = ph
	case p.Ping != nil:
	default:
		panic("probe unspecified")
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			rp.ReportState(probe.State())
		}
	}
}
