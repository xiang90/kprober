package main

import (
	"flag"
	"math/rand"
	"net/http"
)

var (
	sla float64
)

func init() {
	flag.Float64Var(&sla, "sla", 0.5, "Expected SLA")
	flag.Parse()
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if rand.Float64() <= sla {
		w.Write([]byte("OK\n"))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
