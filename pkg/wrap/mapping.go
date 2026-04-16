package wrap

import (
	"regexp"
)

// categoryMapping defines a rule to map an error pattern to a Heimdall category.
type categoryMapping struct {
	patterns []string
	category string
	level    string
}

type compiledRule struct {
	res      []*regexp.Regexp
	category string
	level    string
}

var compiledMappings []compiledRule

// masterMappings is the Master Mapping Table defined in Epic 3.
var masterMappings = []categoryMapping{
	{[]string{"sqlalchemy.exc.OperationalError", "aiomysql.OperationalError", "psycopg2.OperationalError", "pymongo.errors.ConnectionFailure", "gorm\\..*"}, "DATABASE_ERROR", "Persistence"},
	{[]string{"sqlalchemy.exc.ProgrammingError", "pymysql.err.ProgrammingError", "django.db.utils.ProgrammingError"}, "DATABASE_ERROR", "Persistence"},
	{[]string{"mysql.connector.errors.InternalError", "psycopg2.errors.DeadlockDetected", "sqlalchemy.exc.InternalError"}, "CONCURRENCY_ERROR", "Persistence"},
	{[]string{"sqlalchemy.exc.NoResultFound", "sql.ErrNoRows", "django.core.exceptions.ObjectDoesNotExist", "redis.Nil"}, "RESOURCE_NOT_FOUND", "Persistence"},
	{[]string{"sqlalchemy.exc.IntegrityError", "pymysql.err.IntegrityError", "fastapi.exceptions.RequestValidationError", "pydantic.ValidationError"}, "VALIDATION_ERROR", "Persistence/Logic"},
	{[]string{"botocore.exceptions.ClientError.*404", "os.ErrNotExist", "google.cloud.exceptions.NotFound"}, "RESOURCE_NOT_FOUND", "Infrastructure"},
	{[]string{"botocore.exceptions.ClientError.*AccessDenied", "PermissionError", "403 Forbidden", "os.ErrPermission"}, "ACCESS_DENIED", "Security"},
	{[]string{"botocore.exceptions.CredentialRetrievalError", "ExpiredSignatureError", "jwt.exceptions.InvalidTokenError", "jose.exceptions.*"}, "AUTH_EXPIRED", "Security"},
	{[]string{"watchtower.CloudWatchLogBatchError", "boto3.exceptions.*", "google.api_core.exceptions.*"}, "CLOUD_SERVICE_ERROR", "Infrastructure"},
	{[]string{"redis.exceptions.ConnectionError", "redis.exceptions.TimeoutError", "redis.exceptions.BusyLoadingError"}, "CACHE_ERROR", "Cache"},
	{[]string{"rq.exceptions.DequeuingError", "celery.exceptions.TimeoutError", "pika.exceptions.AMQPConnectionError", "github.com/hibiken/asynq.*"}, "TASK_QUEUE_ERROR", "Infrastructure"},
	{[]string{"httpx.HTTPStatusError.*429", "ratelimit.exception.*"}, "RATE_LIMIT_EXCEEDED", "Traffic"},
	{[]string{"httpx.HTTPStatusError.*5[0-9]{2}", "requests.exceptions.HTTPError.*5[0-9]{2}", "pywa.errors.WhatsAppError"}, "UPSTREAM_FAILURE", "Integration"},
	{[]string{"httpx.TimeoutException", "requests.exceptions.Timeout", "context.DeadlineExceeded"}, "INTEGRATION_TIMEOUT", "Upstream"},
	{[]string{"net.OpError", "net.AddrError", "syscall.ECONNREFUSED", "syscall.ETIMEDOUT"}, "NETWORK_ERROR", "System/OS"},
	{[]string{"json.JSONDecodeError", "yaml.YAMLError", "marshmallow.*", "weasyprint.WeasyPrintException", "pandas.errors.*", "numpy.linalg.LinAlgError"}, "DATA_MALFORMED", "Logic"},
	{[]string{"ZeroDivisionError", "KeyError", "AttributeError", "RecursionError", "Panic"}, "LOGIC_CRASH", "System"},
	{[]string{"FileNotFoundError", "IOError", "os.ErrClosed"}, "F_SYSTEM_ERROR", "System"},
}

func init() {
	for _, m := range masterMappings {
		var res []*regexp.Regexp
		for _, p := range m.patterns {
			res = append(res, regexp.MustCompile(p))
		}
		compiledMappings = append(compiledMappings, compiledRule{res, m.category, m.level})
	}
}

// matchPattern checks if an error string matches any of the given patterns using regex objects.
func matchPattern(errStr string, res []*regexp.Regexp) bool {
	for _, re := range res {
		if re.MatchString(errStr) {
			return true
		}
	}
	return false
}

// DetermineCategory infers the Heimdall category and logic level.
func DetermineCategory(errStr string) (string, string) {
	for _, rule := range compiledMappings {
		if matchPattern(errStr, rule.res) {
			return rule.category, rule.level
		}
	}
	return "SYSTEM_ERROR", "System"
}
