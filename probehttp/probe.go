package probehttp

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/xiang90/kprober/pkg/spec"
	"github.com/xiang90/kprober/reporting"
)

var (
	defaultHTTPCheckInterval = time.Second
)

type Probe struct {
	*spec.HTTPProbe
	URL string

	state  reporting.State
	reason string
}

func (p *Probe) Start(ctx context.Context) {
	c := &http.Client{}

	if p.Interval == 0 {
		p.Interval = defaultHTTPCheckInterval
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(p.Interval):
		}

		cctx, cancel := context.WithTimeout(ctx, time.Duration(p.TimeoutSeconds)*time.Second)
		// todo: support body
		r, err := http.NewRequest(p.Method, p.URL, nil)
		if err != nil {
			panic("cannot create http request")
		}
		r.WithContext(cctx)

		resp, err := c.Do(r)
		cancel()
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

	p.checkBody(r.Body)

	p.state = reporting.StateHealthy
	p.reason = ""
}

func (p *Probe) State() (reporting.State, string) {
	return p.state, p.reason
}

func (p *Probe) checkBody(reader io.Reader) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		p.state = reporting.StateDegraded
		p.reason = "failed to read response body"
		return
	}

	match := false
	for _, ep := range p.BodyMatchesRegexp {
		// TODO: cache compile result
		re, err := regexp.Compile(ep)
		if err != nil {
			// TODO: validate regexp before checking!
			continue
		}
		if re.Match(body) {
			match = true
			break
		}
	}
	if len(p.BodyMatchesRegexp) != 0 && !match {
		p.state = reporting.StateDegraded
		p.reason = "body does not match any given regexp"
		return
	}

	for _, expression := range p.BodyDoesNotMatchRegexp {
		// TODO: cache compile result
		re, err := regexp.Compile(expression)
		if err != nil {
			// TODO: validate regexp before checking!
			continue
		}
		if re.Match(body) {
			p.state = reporting.StateDegraded
			p.reason = "body matches a given negative regexp"
			break
		}
	}
}
