package probehttp

import (
	"context"
	"net/http"
	"time"
)

type Probe struct {
	Method   string
	URL      string
	Interval time.Duration

	StatusCode int

	state  int
	Reason string
}

func (p *Probe) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval):
		}

		resp, err := http.Get(p.URL)
		if err != nil {
			p.state = -1
			p.Reason = err.Error()
			continue
		}

		p.check(resp)
		resp.Body.Close()
	}
}

func (p *Probe) check(r *http.Response) {
	if p.StatusCode != 0 && p.StatusCode != r.StatusCode {
		p.state = -1
		p.Reason = "Status code mismatch"
		return
	}

	// check more
}

func (p *Probe) State() int {
	return p.state
}
