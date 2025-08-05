package main

import (
	"log"
	"net/http"
	"service-a/handlers"
	"service-a/tracing"
)

func main() {
	// Initialize tracing
	tracing.InitTracer()

	http.HandleFunc("/cep", handlers.CepHandler)

	log.Println("Service A running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
