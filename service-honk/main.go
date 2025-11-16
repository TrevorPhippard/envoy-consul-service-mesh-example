package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func honkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Honk Honk!!"}`)
}

func main() {
	http.HandleFunc("/honk", honkHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting service-honk on :8080")
	http.ListenAndServe(":8080", nil)
}
