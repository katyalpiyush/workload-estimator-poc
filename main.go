package main

import (
	"workload-estimator-poc/calculator"
	"workload-estimator-poc/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func estimateHandler(w http.ResponseWriter, r *http.Request) {
	var request models.ComputeRequest

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the calculator
	response := calculator.EstimateResources(request)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Create a new router
	router := mux.NewRouter()

	// Define API route
	router.HandleFunc("/estimate", estimateHandler).Methods("POST")

	// Start server
	port := 8080
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
