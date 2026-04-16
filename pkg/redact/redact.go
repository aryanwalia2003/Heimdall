package redact

import (
	"regexp"
)

var (
	// uriRegex matches URI schemes with passwords: scheme://user:password@host
	// We use greedy matching for the middle part to handle passwords that might contain '@'
	uriRegex = regexp.MustCompile(`(?i)([a-z0-9+.-]+://[^:]+:)(.*)(@[^/]+)`)

	// keyRegex matches common sensitive keys in JSON or key=value format (case-insensitive)
	// Supports: key=value, key: value, "key": "value", 'key': 'value'
	keyRegex = regexp.MustCompile(`(?i)("?'?(?:password|secret|api_key|auth_token|jwt|database_url|access_key|token)"?'?[\s]*[:=][\s]*["']?)([^"'\s,]+)(["']?)`)
)

// SanitizeString redacts sensitive terms from the input.
func SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// 1. Redact URI passwords
	res := uriRegex.ReplaceAllString(input, `${1}[REDACTED]${3}`)

	// 2. Redact key-value pairs (case-insensitive)
	res = keyRegex.ReplaceAllString(res, `${1}[REDACTED]${3}`)

	return res
}
