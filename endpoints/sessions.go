package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SessionInfo contains information about a session
// swagger:model SessionInfo
type SessionInfo struct {
	// Original filename of uploaded env file
	// example: production.env
	OriginalFilename string `json:"original_filename"`

	// Upload time
	// example: 2025-08-26T10:30:00Z
	UploadTime string `json:"upload_time"`

	// Number of variables in the session
	// example: 15
	VariablesCount int `json:"variables_count"`
}

// SessionsResponse represents the response for listing sessions
// swagger:model SessionsResponse
type SessionsResponse struct {
	// Response status
	// example: success
	Status string `json:"status"`

	// Number of active sessions
	// example: 3
	ActiveSessions int `json:"active_sessions"`

	// Map of session ID to session info
	Sessions map[string]SessionInfo `json:"sessions"`
}

func GetSessions(apiUrl string) (*SessionsResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl+"/sessions", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var sessions SessionsResponse
	if err := json.Unmarshal(body, &sessions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &sessions, nil
}
