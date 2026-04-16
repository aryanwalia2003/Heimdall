package wrap

import (
	"testing"
)

func TestDetermineCategory(t *testing.T) {
	tests := []struct {
		name          string
		errorPattern  string
		expectedCat   string
		expectedLevel string
	}{
		{"SQLAlchemy Operational", "sqlalchemy.exc.OperationalError", "DATABASE_ERROR", "Persistence"},
		{"MySQL Programming", "pymysql.err.ProgrammingError", "DATABASE_ERROR", "Persistence"},
		{"Deadlock MySQL", "mysql.connector.errors.InternalError", "CONCURRENCY_ERROR", "Persistence"},
		{"Pydantic Validation", "pydantic.ValidationError", "VALIDATION_ERROR", "Persistence/Logic"},
		{"Boto3 Not Found", "botocore.exceptions.ClientError (NoSuchKey/404)", "RESOURCE_NOT_FOUND", "Infrastructure"},
		{"OS Permission", "os.ErrPermission", "ACCESS_DENIED", "Security"},
		{"Expired JWT", "jwt.exceptions.InvalidTokenError", "AUTH_EXPIRED", "Security"},
		{"Boto3 General", "boto3.exceptions.S3UploadFailedError", "CLOUD_SERVICE_ERROR", "Infrastructure"},
		{"Redis Timeout", "redis.exceptions.TimeoutError", "CACHE_ERROR", "Cache"},
		{"Celery Timeout", "celery.exceptions.TimeoutError", "TASK_QUEUE_ERROR", "Infrastructure"},
		{"Rate Limit", "httpx.HTTPStatusError: 429 Too Many Requests", "RATE_LIMIT_EXCEEDED", "Traffic"},
		{"Upstream 5xx", "requests.exceptions.HTTPError: 500 Internal Server Error", "UPSTREAM_FAILURE", "Integration"},
		{"Upstream 5xx Real", "requests.exceptions.HTTPError: 502 Bad Gateway", "UPSTREAM_FAILURE", "Integration"},
		{"Gorm wildcard", "gorm.ErrRecordNotFound", "DATABASE_ERROR", "Persistence"},
		{"Integration Timeout", "requests.exceptions.Timeout", "INTEGRATION_TIMEOUT", "Upstream"},
		{"Network Refused", "syscall.ECONNREFUSED", "NETWORK_ERROR", "System/OS"},
		{"JSON Malformed", "json.JSONDecodeError", "DATA_MALFORMED", "Logic"},
		{"Zero Division", "ZeroDivisionError", "LOGIC_CRASH", "System"},
		{"File Not Found", "FileNotFoundError", "F_SYSTEM_ERROR", "System"},
		{"Unknown Default", "SomeWeirdInternalErrorThatMatchesNothing", "SYSTEM_ERROR", "System"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cat, level := DetermineCategory(tc.errorPattern)
			if cat != tc.expectedCat {
				t.Fatalf("Expected category %s, got %s for %s", tc.expectedCat, cat, tc.errorPattern)
			}
			if level != tc.expectedLevel {
				t.Fatalf("Expected level %s, got %s for %s", tc.expectedLevel, level, tc.errorPattern)
			}
		})
	}
}
