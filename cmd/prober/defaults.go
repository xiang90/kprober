package main

import (
	"time"

	"github.com/xiang90/kprober/pkg/spec"
)

var defaultSpec = spec.ProberSpec{
	Target: spec.Target{
		IP: "www.google.com",
	},
	Probe: spec.Probe{
		Ping: &spec.PingProbe{
			Interval: 1 * time.Second,
			Timeout:  1 * time.Second,
		},
	},
}
