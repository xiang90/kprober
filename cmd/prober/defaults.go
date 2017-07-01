package main

import (
	"github.com/xiang90/kprober/pkg/spec"
)

var defaultSpec = spec.ProberSpec{
	Target: spec.Target{
		IP: "www.google.com",
	},
	Probe: spec.Probe{
		Ping: &spec.PingProbe{
			PeriodSeconds:  1,
			TimeoutSeconds: 1,
		},
	},
}
