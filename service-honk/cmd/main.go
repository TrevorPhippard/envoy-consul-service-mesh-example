package main

import (
	"fmt"
	"net/http"
	"os"
	"service-honk/internal/consul"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func honkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Honk Honk!!"}`)
}

func main() {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	agent := consul.NewAgent(&api.Config{Address: addr})

	serviceCfg := consul.Config{
		ServiceID:   "service-honk-1",
		ServiceName: "service-honk",
		Address:     "service-honk",
		Port:        8080,
		Tags:        []string{"honk"},
		TTL:         8 * time.Second,
		CheckID:     "check_health",
	}

	agent.RegisterService(serviceCfg)

	http.HandleFunc("/honk", honkHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting service-honk on :8080")
	http.ListenAndServe(":8080", nil)
}
