package main

import (
	"fmt"
	"net/http"
	"service-hello/internal/consul"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Hello, world!"}`)
}

func main() {
	agent := consul.NewAgent(&api.Config{Address: "127.0.0.1:8500"})

	serviceCfg := consul.Config{
		ServiceID:   "service-hello-1",
		ServiceName: "service-hello",
		Address:     "localhost",
		Port:        8080,
		Tags:        []string{"hello"},
		TTL:         8 * time.Second,
		CheckID:     "check_health",
	}

	agent.RegisterService(serviceCfg)

	http.HandleFunc("/hello", helloHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting service-hello on :8080")
	http.ListenAndServe(":8080", nil)
}
