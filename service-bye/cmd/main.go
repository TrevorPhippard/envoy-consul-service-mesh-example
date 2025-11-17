package main

import (
	"fmt"
	"log"
	"net/http"
	"service-bye/internal/consul"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func byeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Goodbye, world!"}`)
}

func main() {
	agent := consul.NewAgent(&api.Config{Address: "127.0.0.1:8500"})

	serviceCfg := consul.Config{
		ServiceID:   "service-bye-1",
		ServiceName: "service-bye",
		Address:     "localhost",
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
