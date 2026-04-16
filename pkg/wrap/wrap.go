package wrap

import (
	"encoding/json"

	"github.com/aryanwalia/heimdall/pkg/redact"
)

// RawError represents the incoming error payload from adapters.
type RawError struct {
	TraceID      string `json:"trace_id"`
	ErrorMessage string `json:"error"`
	StackTrace   string `json:"stack_trace"`
	StatusCode   int    `json:"status_code"`
}

// PremiumError represents the generic output.
type PremiumError struct {
	TraceID     string `json:"trace_id"`
	Category    string `json:"category"`
	Message     string `json:"message"`
	IsSanitized bool   `json:"is_sanitized"`
}

// ProcessError parses raw JSON from an Adapter and returns Premium JSON.
func ProcessError(input string) string {
	var raw RawError
	if err := json.Unmarshal([]byte(input), &raw); err != nil {
		return toJSON(PremiumError{Category: "SYSTEM_ERROR", IsSanitized: false})
	}

	cat, _ := DetermineCategory(raw.ErrorMessage)

	// Parse inputs and sanitize
	sanitizedMsg := redact.SanitizeString(raw.ErrorMessage)

	return toJSON(PremiumError{
		TraceID:     raw.TraceID,
		Category:    cat,
		Message:     sanitizedMsg,
		IsSanitized: true,
	})
}

// toJSON serializes struct to JSON, ignoring errors for minimal impl.
func toJSON(resp PremiumError) string {
	b, _ := json.Marshal(resp)
	return string(b)
}
