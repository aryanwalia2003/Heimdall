package wrap

import (
	"encoding/json"
	"testing"
)

type ErrorResponse struct {
	TraceID     string `json:"trace_id"`
	Category    string `json:"category"`
	Message     string `json:"message"`
	IsSanitized bool   `json:"is_sanitized"`
}

func TestProcessError_InvalidJSON(t *testing.T) {
	input := `invalid_json`
	output := ProcessError(input)

	var resp ErrorResponse
	err := json.Unmarshal([]byte(output), &resp)
	if err != nil {
		t.Fatalf("ProcessError did not return valid JSON for invalid input: %v", err)
	}

	if resp.Category != "SYSTEM_ERROR" {
		t.Errorf("Expected category SYSTEM_ERROR, got %v", resp.Category)
	}
}

func TestProcessError_ValidBasicJSON(t *testing.T) {
	input := `{"error_message": "test error", "stack_trace": "test stack", "status_code": 500}`
	output := ProcessError(input)

	var resp ErrorResponse
	err := json.Unmarshal([]byte(output), &resp)
	if err != nil {
		t.Fatalf("ProcessError did not return valid JSON: %v", err)
	}

	if resp.IsSanitized != true {
		t.Errorf("Expected IsSanitized true, got false")
	}
}
