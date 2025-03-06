package main

import (
	"assignment_1/handlers"
	"log"
	"net/http"
	"time"
)

func main() {
	PORT := ":8080"

	router := http.NewServeMux()

	handlers.SetStartTime(time.Now())

	// Default handler endpoint that points to the other endpoints
	router.HandleFunc(handlers.DEFAULT_PATH, handlers.DefaultHandler)

	// Info handler endpoint
	router.HandleFunc(handlers.INFO_PATH, handlers.InfoHandler)

	// Population handler endpoint
	router.HandleFunc(handlers.POPULATION_PATH, handlers.PopulationHandler)

	// Status handler endpoint
	router.HandleFunc(handlers.STATUS_PATH, handlers.StatusHandler)

	log.Println("Starting server on port " + PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
