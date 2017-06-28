package probehttp

import (
	"context"
	"net/http"
	"time"

	"github.com/xiang90/kprober/reporting"
)

var (
	defaultHTTPCheckInterval = time.Second
)

type Probe struct {
	Method   string
	URL      string
	Interval time.Duration

	StatusCode int

	state  reporting.State
	reason string
}

func (p *Probe) Start(ctx context.Context) {
	if p.Interval == 0 {
		p.Interval = defaultHTTPCheckInterval
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval):
		}

		resp, err := http.Get(p.URL)
		if err != nil {
			p.state = reporting.StateDown
			p.reason = err.Error()
			continue
		}

		p.check(resp)
		resp.Body.Close()
	}
}

func (p *Probe) check(r *http.Response) {
	if p.StatusCode != 0 && p.StatusCode != r.StatusCode {
		p.state = reporting.StateDown
		p.reason = "Status code mismatch"
		return
	}

	// check more

	p.state = reporting.StateHealthy
	p.reason = ""
}

func (p *Probe) State() (reporting.State, string) {
	return p.state, p.reason
}
