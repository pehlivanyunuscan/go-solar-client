package endpoints

import (
	"fmt"
	"net/http"
)

func DeleteSession(apiUrl, sessionID string) error {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", apiUrl+"/sessions/"+sessionID, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("session not found: %s", sessionID)
	case http.StatusInternalServerError:
		return fmt.Errorf("server error while deleting session: %s", sessionID)
	default:
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
