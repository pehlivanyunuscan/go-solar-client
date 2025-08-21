package main

import (
	"encoding/json"
	"fmt"
	"go-solar-client/endpoints"
	"io"
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

	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}
		// Parse the request body into a RunRequest struct
		var req endpoints.RunRequest
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}
		resp, err := endpoints.RunForecast(apiUrl, &req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to run forecast: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Starting server on :8888")
	http.ListenAndServe(":8888", nil)
}
