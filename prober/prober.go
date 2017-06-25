package prober

import (
	"context"
	"fmt"
	"time"

	"github.com/xiang90/kprober/k8sutil"
	"github.com/xiang90/kprober/probehttp"
	"github.com/xiang90/kprober/probeping"
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
	if p.Namespace == "" {
		p.Namespace = "default"
	}

	var (
		ip  string
		err error
	)
	switch {
	case p.Pod != nil:
		ip, err = k8sutil.IPFromPod(p.Namespace, *p.Pod)
	default:
		panic("target unspecified")
	}

	rp := reporting.NewPrometheus(p.Name)

	var probe Probe

	switch {
	case p.HTTP != nil:
		url := fmt.Sprintf("%s://%s:%s/%s", p.HTTP.Scheme, ip, p.HTTP.Port, p.HTTP.Path)
		ph := &probehttp.Probe{
			URL:      url,
			Method:   p.HTTP.Method,
			Interval: p.HTTP.Interval,
		}
		ph.Start(context.TODO())
		probe = ph
	case p.Ping != nil:
		pp := &probeping.Probe{
			Addr:     ip,
			Interval: p.HTTP.Interval,
		}
		pp.Start(context.TODO())
		probe = pp
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
