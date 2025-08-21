package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RunRequest represents the structure of the request body for running a forecast.
// It contains fields that are required to perform the forecast operation.
type RunRequest struct {
	PrometheusURL       string  `json:"PROMETHEUS_URL"`
	MetricName          string  `json:"METRIC_NAME"`
	TrainDays           int     `json:"TRAIN_DAYS"`
	BatteryCapacityWh   float64 `json:"BATTERY_CAPACITY_WH"`
	InitialSocPercent   float64 `json:"INITIAL_SOC_PERCENT"`
	ConstantLoadW       float64 `json:"CONSTANT_LOAD_W"`
	ChargeEfficiency    float64 `json:"CHARGE_EFFICIENCY"`
	DischargeEfficiency float64 `json:"DISCHARGE_EFFICIENCY"`
	DetailedSummary     bool    `json:"DETAILED_SUMMARY"`
	UseCython           bool    `json:"USE_CYTHON"`
}

type RunResponse struct {
	Status    string                 `json:"status"`
	Result    map[string]interface{} `json:"result,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

func RunForecast(apiUrl string, reqData *RunRequest) (*RunResponse, error) {
	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(apiUrl+"/run", "application/json", bytes.NewBuffer(body))
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

	var runResp RunResponse
	if err := json.Unmarshal(respBody, &runResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &runResp, nil
}
