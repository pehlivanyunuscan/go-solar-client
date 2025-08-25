package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SessionInfo struct {
	OriginalFilename string `json:"original_filename"`
	UploadTime       string `json:"upload_time"`
	VariablesCount   int    `json:"variables_count"`
}

type SessionsResponse struct {
	ActiveSessions int                    `json:"active_sessions"`
	Sessions       map[string]SessionInfo `json:"sessions"`
	Status         string                 `json:"status"`
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
