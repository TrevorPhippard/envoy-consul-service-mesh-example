package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"service-bye/internal/consul"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func byeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Goodbye, world!"}`)
}

func main() {
	addr := os.Getenv("CONSUL_HTTP_ADDR")
	agent := consul.NewAgent(&api.Config{Address: addr})

	serviceCfg := consul.Config{
		ServiceID:   "service-bye-1",
		ServiceName: "service-bye",
		Address:     "service-bye",
		Port:        8080,
		Tags:        []string{"bye"},
		TTL:         8 * time.Second,
		CheckID:     "check_health",
	}

	agent.RegisterService(serviceCfg)

	http.HandleFunc("/bye", byeHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("service-bye running on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
