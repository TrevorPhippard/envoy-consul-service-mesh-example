package main

import (
	"fmt"
	"net/http"
	"os"
	"service-hello/internal/consul"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Hello, world!"}`)
}

func main() {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	agent := consul.NewAgent(&api.Config{Address: addr})

	serviceCfg := consul.Config{
		ServiceID:   "service-hello-1",
		ServiceName: "service-hello",
		Address:     "service-hello",
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
