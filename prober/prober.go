package prober

import (
	"context"
	"fmt"
	"time"

	"github.com/xiang90/kprober/pkg/spec"
	"github.com/xiang90/kprober/pkg/util/k8sutil"
	"github.com/xiang90/kprober/probehttp"
	"github.com/xiang90/kprober/probeping"
	"github.com/xiang90/kprober/reporting"
)

type Probe interface {
	State() (reporting.State, string)
}

type Prober struct {
	Name string
	spec.ProberSpec
}

func (p *Prober) Start(ctx context.Context) {
	ts := p.ProberSpec.Target
	ps := p.ProberSpec.Probe

	var (
		ip  string
		err error
	)
	switch {
	case ts.Pod != nil:
		pod := ts.Pod
		ip, err = k8sutil.IPFromPod(pod.Namespace, pod.Name)
		if err != nil {
			// TODO: retry and report pod as unhealthy
			panic(err)
		}
	case ts.IP != "":
		ip = ts.IP
	default:
		panic("target unspecified")
	}

	rp := reporting.NewPrometheus(p.Name)

	var probe Probe

	switch {
	case ps.HTTP != nil:
		url := fmt.Sprintf("%s://%s:%s/%s", ps.HTTP.Scheme, ip, ps.HTTP.Port, ps.HTTP.Path)
		ph := &probehttp.Probe{
			URL:      url,
			Method:   ps.HTTP.Method,
			Interval: ps.HTTP.Interval,

			StatusCode: ps.HTTP.StatusCode,
		}
		go ph.Start(context.TODO())
		probe = ph
	case ps.Ping != nil:
		pp := &probeping.Probe{
			Addr:     ip,
			Interval: ps.Ping.Interval,
		}
		go pp.Start(context.TODO())
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
