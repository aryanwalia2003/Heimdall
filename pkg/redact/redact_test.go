package redact

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No sensitive info",
			input:    "This is a normal error log",
			expected: "This is a normal error log",
		},
		{
			name:     "Contains password",
			input:    "Error connecting, password=supersecret, retry",
			expected: "Error connecting, password=[REDACTED], retry",
		},
		{
			name:     "Contains jwt",
			input:    `{"error": "invalid token", "jwt": "eyJhb..."}`,
			expected: `{"error": "invalid token", "jwt": "[REDACTED]"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SanitizeString(tc.input)
			if result != tc.expected {
				t.Errorf("SanitizeString failed\nExpected: %s\nGot:      %s", tc.expected, result)
			}
		})
	}
}
