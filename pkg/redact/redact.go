package redact

import "strings"

var sensitiveKeys = []string{
	"password", "secret", "api_key", "auth_token", "jwt", "database_url", "access_key",
}

// SanitizeString redacts sensitive terms from the input.
func SanitizeString(input string) string {
	res := input
	for _, key := range sensitiveKeys {
		res = replaceAssign(res, key+"=")
		res = replaceJSON(res, `"`+key+`": "`)
	}
	return res
}

func replaceAssign(input, prefix string) string {
	var sb strings.Builder
	for {
		idx := strings.Index(input, prefix)
		if idx == -1 {
			sb.WriteString(input)
			break
		}
		sb.WriteString(input[:idx+len(prefix)])
		sb.WriteString("[REDACTED]")

		rest := input[idx+len(prefix):]
		end := strings.IndexAny(rest, " \n\r\t,")
		if end == -1 {
			break
		}
		input = rest[end:]
	}
	return sb.String()
}

func replaceJSON(input, prefix string) string {
	var sb strings.Builder
	for {
		idx := strings.Index(input, prefix)
		if idx == -1 {
			sb.WriteString(input)
			break
		}
		sb.WriteString(input[:idx+len(prefix)])
		sb.WriteString("[REDACTED]")

		rest := input[idx+len(prefix):]
		end := strings.IndexByte(rest, '"')
		if end == -1 {
			break
		}
		input = rest[end:]
	}
	return sb.String()
}
