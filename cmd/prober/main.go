package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/xiang90/kprober/prober"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsListenAddr string
	specFile          string
	name              string
)

func init() {
	flag.StringVar(&metricsListenAddr, "metrics-listen-addr", "0.0.0.0:17783", "listen address")
	flag.StringVar(&specFile, "f", "spec.yaml", "prober spec file in yaml format")
	flag.StringVar(&name, "n", "default-prober", "name of the prober")
	flag.Parse()
}

func main() {

	go func() {
		p := prober.Prober{
			Name:       name,
			ProberSpec: defaultSpec,
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
