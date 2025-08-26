package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RunRequest represents the structure of the request body for running a forecast.
// swagger:model RunRequest
type RunRequest struct {
	// The Prometheus URL to fetch data from
	// required: true
	// example: http://localhost:9090
	PrometheusURL string `json:"PROMETHEUS_URL"`

	// The metric name to query from Prometheus
	// required: true
	// example: solar_power
	MetricName string `json:"METRIC_NAME"`

	// Number of days to use for training
	// required: true
	// example: 7
	// minimum: 1
	TrainDays int `json:"TRAIN_DAYS"`

	// Battery capacity in Wh
	// example: 5000.0
	BatteryCapacityWh float64 `json:"BATTERY_CAPACITY_WH"`

	// Initial SOC percentage
	// example: 80.0
	InitialSocPercent float64 `json:"INITIAL_SOC_PERCENT"`

	// Constant load in Watts
	// example: 100.0
	ConstantLoadW float64 `json:"CONSTANT_LOAD_W"`

	// Charge efficiency
	// example: 0.9
	ChargeEfficiency float64 `json:"CHARGE_EFFICIENCY"`

	// Discharge efficiency
	// example: 0.9
	DischargeEfficiency float64 `json:"DISCHARGE_EFFICIENCY"`

	// Whether to include detailed summary
	// example: true
	DetailedSummary bool `json:"DETAILED_SUMMARY"`

	// Whether to use Cython
	// example: true
	UseCython bool `json:"USE_CYTHON"`
}

// RunResponse represents the response from running a forecast
// swagger:model RunResponse
type RunResponse struct {
	// Status of the operation
	// example: success
	Status string `json:"status"`

	// Result data
	Result map[string]interface{} `json:"result,omitempty"`

	// Message from the operation
	// example: Forecast completed successfully
	Message string `json:"message,omitempty"`

	// Timestamp of the operation
	// example: 2025-08-26T10:30:00Z
	Timestamp string `json:"timestamp,omitempty"`
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
