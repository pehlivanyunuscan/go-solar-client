package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HealthResponse represents the structure of the health check response.
// swagger:model HealthResponse
type HealthResponse struct {
	// Service name
	// example: Solar Forecasting API
	Service string `json:"service"`

	// Service status
	// example: healthy
	Status string `json:"status"`

	// Health check timestamp
	// example: 2025-08-26T10:30:00Z
	Timestamp string `json:"timestamp,omitempty"`
}

func CheckHealth(apiUrl string) (*HealthResponse, error) {
	client := &http.Client{}
	// Create a new HTTP GET request to the health endpoint
	req, err := http.NewRequest("GET", apiUrl+"/health", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the request header to accept JSON responses
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK (200)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	// Read the response body and unmarshal it into a HealthResponse struct
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	// Unmarshal the JSON response into a HealthResponse struct
	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &health, nil
}
