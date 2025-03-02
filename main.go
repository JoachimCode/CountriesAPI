package main

import (
	"assignment_1/handlers"
	"log"
	"net/http"
)

func main() {
	PORT := ":8080"

	router := http.NewServeMux()

	// Default handler endpoint that points to the other endpoints
	router.HandleFunc(handlers.DEFAULT_PATH, handlers.DefaultHandler)

	// Info handler endpoint
	router.HandleFunc(handlers.INFO_PATH, handlers.InfoHandler)

	log.Println("Starting server on port " + PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
