package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RunWithEnvResponse struct {
	Result    map[string]interface{} `json:"result"`
	SessionID string                 `json:"session_id"`
	Status    string                 `json:"status"`
	Timestamp string                 `json:"timestamp"`
}

func RunWithEnv(apiUrl, sessionID string, override map[string]interface{}) (*RunWithEnvResponse, error) {
	var jsonData []byte
	var err error

	if len(override) > 0 {
		jsonData, err = json.Marshal(override)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal override: %w", err)
		}
	} else {
		jsonData = []byte("{}")
	}

	endpoint := fmt.Sprintf("%s/run-with-env/%s", apiUrl, sessionID)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var result RunWithEnvResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}
