// Package classification Solar Power Forecasting API
//
// # Documentation for Solar Power Forecasting API
//
// The purpose of this application is to provide a REST API for solar power forecasting
//
//	Schemes: http
//	Host: localhost:8888
//	BasePath: /
//	Version: 1.0.0
//	License: MIT https://opensource.org/licenses/MIT
//	Contact: Solar API Team<support@solar-api.com>
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
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

	// swagger:operation GET /health health healthCheck
	// ---
	// summary: Check service health
	// description: Returns the health status of the forecasting service
	// tags:
	// - health
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Service is healthy
	//     schema:
	//       $ref: "#/definitions/HealthResponse"
	//   '500':
	//     description: Service is unhealthy
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// swagger:operation POST /run forecast runForecast
	// ---
	// summary: Run solar power forecast
	// description: Run a solar power forecast with the provided parameters
	// tags:
	// - forecast
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   description: Forecast parameters
	//   required: true
	//   schema:
	//     $ref: "#/definitions/RunRequest"
	// responses:
	//   '200':
	//     description: Forecast completed successfully
	//     schema:
	//       $ref: "#/definitions/RunResponse"
	//   '400':
	//     description: Bad request
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// swagger:operation POST /upload-env environment uploadEnv
	// ---
	// summary: Upload environment file
	// description: Upload a .env file to create a session
	// tags:
	// - environment
	// produces:
	// - application/json
	// parameters:
	// - name: envfile
	//   in: query
	//   description: Path to the environment file
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: File uploaded successfully
	//     schema:
	//       $ref: "#/definitions/UploadEnvResponse"
	//   '400':
	//     description: Bad request
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// swagger:operation POST /run-with-env/{sessionID} forecast runWithEnv
	// ---
	// summary: Run forecast with environment session
	// description: Run forecast using a previously uploaded environment session
	// tags:
	// - forecast
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: sessionID
	//   in: path
	//   description: Session ID from uploaded environment
	//   required: true
	//   type: string
	// - name: body
	//   in: body
	//   description: Optional parameter overrides
	//   schema:
	//     type: object
	// responses:
	//   '200':
	//     description: Forecast completed successfully
	//     schema:
	//       $ref: "#/definitions/RunWithEnvResponse"
	//   '400':
	//     description: Bad request
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
	http.HandleFunc("/run-with-env/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get the session ID from the URL path
		sessionID := r.URL.Path[len("/run-with-env/"):]
		if sessionID == "" {
			http.Error(w, "session_id is required", http.StatusBadRequest)
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

	// swagger:operation GET /sample-env environment getSampleEnv
	// ---
	// summary: Get sample environment
	// description: Returns a sample environment configuration
	// tags:
	// - environment
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Sample environment configuration
	//     schema:
	//       $ref: "#/definitions/SampleEnvResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// swagger:operation GET /sessions sessions listSessions
	// ---
	// summary: List all sessions
	// description: Returns a list of all active environment sessions
	// tags:
	// - sessions
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: List of sessions
	//     schema:
	//       $ref: "#/definitions/SessionsResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// swagger:operation DELETE /sessions/{sessionID} sessions deleteSession
	// ---
	// summary: Delete a session
	// description: Delete a specific environment session by ID
	// tags:
	// - sessions
	// produces:
	// - application/json
	// parameters:
	// - name: sessionID
	//   in: path
	//   description: Session ID to delete
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Session deleted successfully
	//     schema:
	//       type: object
	//       properties:
	//         message:
	//           type: string
	//         status:
	//           type: string
	//   '400':
	//     description: Bad request
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
	//   '500':
	//     description: Internal server error
	//     schema:
	//       $ref: "#/definitions/ErrorResponse"
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

	// Swagger UI endpoints
	http.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "./swagger.json")
	})

	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Solar Power Forecasting API - Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
        .swagger-ui .topbar { display: none; }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
<script>
window.onload = function() {
    const ui = SwaggerUIBundle({
        url: '/swagger.json',
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        validatorUrl: null
    });
};
</script>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	// Root redirect
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/docs", http.StatusMovedPermanently)
			return
		}
		http.NotFound(w, r)
	})

	fmt.Println("🚀 Starting Solar Power Forecasting API server...")
	fmt.Println("📍 Server: http://localhost:8888")
	fmt.Println("📖 API Documentation: http://localhost:8888/docs")
	fmt.Println("📄 Swagger JSON: http://localhost:8888/swagger.json")
	fmt.Println("💚 Health Check: http://localhost:8888/health")

	http.ListenAndServe(":8888", nil)
}

// Swagger model definitions for responses

// swagger:model ErrorResponse
type ErrorResponse struct {
	// Error message
	// example: An error occurred
	Error string `json:"error"`

	// Error timestamp
	// example: 2025-08-26T10:30:00Z
	Timestamp string `json:"timestamp,omitempty"`
}
