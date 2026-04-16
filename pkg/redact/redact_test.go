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
		{
			name:     "Contains postgres URI",
			input:    "Critical Failure: Connection to postgresql://admin:P@ssword123@zorms-db:5432/dev failed",
			expected: "Critical Failure: Connection to postgresql://admin:[REDACTED]@zorms-db:5432/dev failed",
		},
		{
			name:     "Contains mixed-case api_key",
			input:    "API_KEY=AKIA_MOCK_123456789",
			expected: "API_KEY=[REDACTED]",
		},
		{
			name:     "Contains nested JSON auth_token",
			input:    `"Auth_Token": "secret_v1"`,
			expected: `"Auth_Token": "[REDACTED]"`,
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
