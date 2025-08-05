package main

import (
	"log"
	"net/http"
	"service-b/handlers"
	"service-b/tracing"
)

func main() {
	// Initialize tracing
	tracing.InitTracer()

	http.HandleFunc("/process", handlers.CepHandler)

	log.Println("Service B running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
