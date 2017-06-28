package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xiang90/kprober/pkg/client"
	"github.com/xiang90/kprober/prober"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsListenAddr string
	specFile          string
	name              string
	namespace         string
)

func init() {
	flag.StringVar(&metricsListenAddr, "metrics-listen-addr", "0.0.0.0:17783", "listen address")
	flag.StringVar(&name, "n", "default-prober", "name of the prober")
	flag.StringVar(&namespace, "ns", "default", "namespace of the prober")

	flag.Parse()
}

func main() {
	spec := defaultSpec

	// todo: init client
	var pc client.ProbersCR

	crd, err := pc.Get(context.TODO(), namespace, name)
	if err != nil {
		fmt.Println(err)
	}
	if crd != nil {
		fmt.Println("using prober spec from CRD")
		spec = crd.Spec
	} else {
		fmt.Println("using default spec")
	}

	go func() {
		p := prober.Prober{
			Name:       name,
			ProberSpec: spec,
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
