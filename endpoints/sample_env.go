package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SampleEnvResponse represents the response for sample environment
// swagger:model SampleEnvResponse
type SampleEnvResponse struct {
	// Response status
	// example: success
	Status string `json:"status"`

	// Description of the sample environment
	// example: Sample environment configuration for solar forecasting
	Description string `json:"description"`

	// Sample environment file content
	SampleEnvContent string `json:"sample_env_content"`
}

func GetSampleEnv(apiUrl string) (*SampleEnvResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", apiUrl+"/sample-env", nil)
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

	var sampleEnvResp SampleEnvResponse
	if err := json.Unmarshal(body, &sampleEnvResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return &sampleEnvResp, nil
}
