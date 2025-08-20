package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HealthResponse represents the structure of the health check response.
type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// CheckHealth sends a health check request to the specified API URL and returns the health status.
// It returns a HealthResponse and an error if any occurs.

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
