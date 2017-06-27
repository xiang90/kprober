package probeping

import (
	"context"
	"time"

	"github.com/sparrc/go-ping"
)

var (
	defaultPingInterval = time.Second
)

type Probe struct {
	Addr     string
	Interval time.Duration

	MaxLatency time.Duration

	state  int
	Reason string
}

func (p *Probe) Start(ctx context.Context) {
	if p.Interval == 0 {
		p.Interval = defaultPingInterval
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval):
		}

		pinger, err := ping.NewPinger(p.Addr)
		if err != nil {
			p.state = -1
			p.Reason = err.Error()
			continue
		}

		pinger.Timeout = time.Second

		pinger.Count = 1
		pinger.Run()
		p.check(pinger.Statistics())
	}
}

func (p *Probe) check(s *ping.Statistics) {
	if s.PacketLoss != 0 {
		p.state = -1
		p.Reason = "ping packet lost"
		return
	}

	// check more
}

func (p *Probe) State() int {
	return p.state
}
