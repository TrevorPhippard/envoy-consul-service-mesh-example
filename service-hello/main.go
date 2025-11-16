package main

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{"message":"Hello, world!"}`)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Starting service-hello on :8080")
	http.ListenAndServe(":8080", nil)
}
