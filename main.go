package main

import (
	"encoding/json"
	"fmt"
	"go-solar-client/endpoints"
	"net/http"
)

func main() {
	apiUrl := "http://10.67.67.25:4545"

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		health, err := endpoints.CheckHealth(apiUrl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Health check failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(health); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})
	fmt.Println("Starting server on :8888")
	http.ListenAndServe(":8888", nil)
}
