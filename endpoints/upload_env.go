package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadEnvResponse response structure for upload env file
// swagger:model UploadEnvResponse
type UploadEnvResponse struct {
	// Response message
	// example: Environment file uploaded successfully
	Message string `json:"message"`

	// Session ID created for this upload
	// example: abc123def456
	SessionID string `json:"session_id"`

	// Status of the upload
	// example: success
	Status string `json:"status"`

	// List of environment variables found
	Variables []string `json:"variables"`

	// Number of variables found
	// example: 10
	VariablesCount int `json:"variables_count"`
}

// UploadEnvFile verilen env dosyasını /upload-env endpointine yükler ve yanıtı döner.
func UploadEnvFile(apiUrl, envFilePath string) (*UploadEnvResponse, error) {
	file, err := os.Open(envFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a buffer to hold the multipart form data
	var buf bytes.Buffer
	// Create a new multipart writer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("env_file", "params.env")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file content: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", apiUrl+"/upload-env", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// Set the correct content type for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
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

	var uploadResp UploadEnvResponse
	if err := json.Unmarshal(respBody, &uploadResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &uploadResp, nil
}
