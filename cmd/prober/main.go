package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/xiang90/kprober/prober"
	"github.com/xiang90/kprober/spec"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsListenAddr string
)

func init() {
	flag.StringVar(&metricsListenAddr, "metrics-listen-addr", "0.0.0.0:17783", "listen address")
	flag.Parse()
}

func main() {

	go func() {
		googleAddr := "www.google.com"
		p := prober.Prober{
			ProbeSpec: spec.ProbeSpec{
				Target: spec.Target{
					IP: &googleAddr,
				},
				Probe: spec.Probe{
					Ping: &spec.PingProbe{
						Interval: 1 * time.Second,
					},
				},
			},
		}

		p.Start(context.TODO())
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Serving metrics on %v", metricsListenAddr)

	s := &http.Server{
		Addr:           metricsListenAddr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}
