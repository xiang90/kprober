package probecontainer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/xiang90/kprober/reporting"
)

type Probe struct {
	ResultReader io.ReadCloser

	state  reporting.State
	reason string
}

func (p *Probe) Start(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			p.ResultReader.Close()
			return
		}
	}()

	br := bufio.NewReader(p.ResultReader)
	for {
		var (
			state  int
			reason string
		)

		l, err := br.ReadString('\n')
		if err != nil {
			return
		}
		n, err := fmt.Sscanf(l, "%d %s", &state, &reason)
		if n != 2 && err != io.EOF {
			p.state = reporting.StateUnknown
			p.reason = "invalid result"
			log.Printf("invalid result from the probe container: %s", l)
			continue
		}

		p.state = reporting.State(n)
		p.reason = reason
	}
}

func (p *Probe) State() (reporting.State, string) {
	return p.state, p.reason
}
