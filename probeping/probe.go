package probeping

import (
	"context"
	"time"

	"github.com/xiang90/kprober/reporting"

	"github.com/sparrc/go-ping"
)

var (
	defaultPingInterval = time.Second
	defaultPingTimeout  = time.Second
)

type Probe struct {
	Addr string

	Interval time.Duration
	Timeout  time.Duration

	state  reporting.State
	reason string
}

func (p *Probe) Start(ctx context.Context) {
	if p.Interval == 0 {
		p.Interval = defaultPingInterval
	}
	if p.Timeout == 0 {
		p.Timeout = defaultPingTimeout
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval):
		}

		pinger, err := ping.NewPinger(p.Addr)
		pinger.SetPrivileged(true)
		if err != nil {
			p.state = reporting.StateDown
			p.reason = err.Error()
			continue
		}

		pinger.Timeout = p.Timeout

		pinger.Count = 1
		pinger.Run()
		p.check(pinger.Statistics())
	}
}

func (p *Probe) check(s *ping.Statistics) {
	if s.PacketLoss != 0 {
		p.state = reporting.StateDown
		p.reason = "ping packet lost"
		return
	}
	p.state = reporting.StateHealthy
	p.reason = ""

	// check more
}

func (p *Probe) State() (reporting.State, string) {
	return p.state, p.reason
}
