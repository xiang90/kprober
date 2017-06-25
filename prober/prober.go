package prober

import (
	"context"
	"time"

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
	rp := reporting.NewPrometheus(p.Name)

	var probe Probe

	switch {
	case p.HTTP != nil:
	case p.Ping != nil:
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
