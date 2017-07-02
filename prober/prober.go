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
	ts, ps := p.ProberSpec.Target, p.ProberSpec.Probe
	kubecli := k8sutil.MustNewKubeClient()

	var (
		ip  string
		err error
	)
	switch {
	case ts.Service != nil:
		srv := ts.Service
		ip, err = k8sutil.IPFromService(kubecli, srv.Namespace, srv.Name)
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
		hp := ps.HTTP
		url := fmt.Sprintf("%s://%s:%d/%s", hp.Scheme, ip, hp.Port, hp.Path)
		ph := &probehttp.Probe{
			URL:       url,
			HTTPProbe: hp,
		}
		go ph.Start(context.TODO())
		probe = ph
	case ps.Ping != nil:
		pp := &probeping.Probe{
			Addr:     ip,
			Interval: time.Duration(ps.Ping.PeriodSeconds) * time.Second,
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
