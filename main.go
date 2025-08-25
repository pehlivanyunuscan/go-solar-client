package main

import (
	"encoding/json"
	"fmt"
	"go-solar-client/endpoints"
	"io"
	"net/http"
	"os"
)

func main() {
	// apiUrl := "http://10.67.67.25:4545"
	apiUrl := "http://localhost:4545"

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

	http.HandleFunc("/upload-env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// .env dosya yolunu query parametresi olarak al
		envFile := r.URL.Query().Get("envfile")
		if envFile == "" {
			http.Error(w, "envfile query parameter is required", http.StatusBadRequest)
			return
		}

		// Dosyanın var olup olmadığını kontrol et
		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			http.Error(w, fmt.Sprintf("File does not exist: %s", envFile), http.StatusBadRequest)
			return
		}

		resp, err := endpoints.UploadEnvFile(apiUrl, envFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to upload .env file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/run-with-env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the session ID from the query parameters
		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "session_id query parameter is required", http.StatusBadRequest)
			return
		}

		// Parse the request body into a map for overrides
		var override map[string]interface{}
		// Read the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
			return
		}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &override); err != nil {
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}
		}
		resp, err := endpoints.RunWithEnv(apiUrl, sessionID, override)
		if err != nil {
			http.Error(w, fmt.Sprintf("RunWithEnv failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/sample-env", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}
		resp, err := endpoints.GetSampleEnv(apiUrl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get sample env: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
			return
		}
		sessions, err := endpoints.GetSessions(apiUrl)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get sessions: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sessions)
	})

	http.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
			return
		}
		sessionID := r.URL.Path[len("/sessions/"):]
		if sessionID == "" {
			http.Error(w, "session_id query parameter is required", http.StatusBadRequest)
			return
		}
		err := endpoints.DeleteSession(apiUrl, sessionID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete session: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Session deleted successfully",
			"status":  "success",
		})
	})

	fmt.Println("Starting server on :8888")
	http.ListenAndServe(":8888", nil)
}
