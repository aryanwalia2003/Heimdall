# 🛡️ Project Heimdall: The Ecosystem Error Orchestration Roadmap

This document provides a **painfully descriptive**, step-by-step technical plan to implement a unified error handling and masking system across the entire ZFW microservice ecosystem (`zorms`, `z-or-cs`, `zfw_dashboard`).

---

## 🏛️ Epic 1: The All-Seeing Eye (Core Go Engine)
**Goal**: Build a single, high-performance Go repository that defines "how an error looks" and "how it is sanitized."

### 📖 1.1 Repo Initialization & The "Go Interface"
- Create `zfw-heimdall-core`.
- Implement a `wrap` package with a public function `ProcessError(input string) string`.
- **Logic**: The function takes a raw JSON string from an Adapter (containing the stack trace, error message, and original status code) and returns the "Premium JSON".

### 🔎 1.2 The Redaction Sub-Engine
- Build a dictionary of sensitive terms: `["password", "secret", "api_key", "auth_token", "jwt", "database_url", "access_key"]`.
- Implement a scanner that traverses the raw stack trace and replaces any string matching these keys (and their values) with `[REDACTED]`.

### 🏗️ 1.3 The C-Shared Build Pipeline
- Use `go build -buildmode=c-shared -o libheimdall.so main.go`.
- This step is critical: it allows any language (Python via `ctypes`, Go via `cgo`) to load the exact same compiled logic into memory at native speed.

---

## 🌈 Epic 2: The Bifrost Connectors (Language-Specific Adapters)
**Goal**: Create lightweight libraries that services import to "talk" to Heimdall.

### 🐍 2.1 Python/FastAPI Adapter (`zfw-heimdall-python`)
**Step-by-step Implementation**:
1.  **Middleware Class**: Create a `HeimdallMiddleware` class inheriting from `BaseHTTPMiddleware`.
2.  **Early Trace Inception**: Inside `dispatch()`, check for `X-Trace-ID` in headers. If missing, generate a new ULID. Store this in a `ContextVar` so it’s available globally in the thread.
3.  **Exception Interception**: Wrap the `call_next(request)` in a massive `try/except` block.
4.  **Premium Wrap**: If an exception occurs:
    - Capture the `traceback.format_exc()`.
    - Extract the category hint (check for `X-Error-Category` header or custom class attribute).
    - Call the `libheimdall.so` native function.
5.  **Fail-Open Safety**: If `libheimdall.so` throws a segmentation fault or is missing, the Python code MUST catch that and return a hardcoded JSON string: `{"success": false, "message": "Critical System Error"}`.

### 🐹 2.2 Go Adapter (`zfw-heimdall-go`)
- Create a `HeimdallHandler` for services like `z-or-cs`.
- Since the core logic is already in Go, the adapter just imports the `wrap` package directly (no `.so` overhead needed for Go apps).

---

## 📚 Epic 3: The Asgardian Registry (Detailed Error Mapping)
**Goal**: Mapping every possible library-specific exception to a Heimdall Category.

### 🗺️ 3.1 The Master Mapping Table
The Go Core and Adapters must follow this mapping strictly:

| Exception Pattern (Regex/Class) | Heimdall Category | Logic Level |
| :--- | :--- | :--- |
| `sqlalchemy.exc.OperationalError`, `aiomysql.OperationalError`, `psycopg2.OperationalError`, `pymongo.errors.ConnectionFailure`, `gorm.*` | `DATABASE_ERROR` | Persistence |
| `sqlalchemy.exc.ProgrammingError`, `pymysql.err.ProgrammingError`, `django.db.utils.ProgrammingError` | `DATABASE_ERROR` | Persistence |
| `mysql.connector.errors.InternalError` (1213), `psycopg2.errors.DeadlockDetected`, `sqlalchemy.exc.InternalError` (Lock wait) | `CONCURRENCY_ERROR` | Persistence |
| `sqlalchemy.exc.NoResultFound`, `sql.ErrNoRows`, `django.core.exceptions.ObjectDoesNotExist`, `redis.Nil` | `RESOURCE_NOT_FOUND` | Persistence |
| `sqlalchemy.exc.IntegrityError`, `pymysql.err.IntegrityError`, `fastapi.exceptions.RequestValidationError`, `pydantic.ValidationError` | `VALIDATION_ERROR` | Persistence/Logic |
| `botocore.exceptions.ClientError` (NoSuchKey/404), `os.ErrNotExist`, `google.cloud.exceptions.NotFound` | `RESOURCE_NOT_FOUND` | Infrastructure |
| `botocore.exceptions.ClientError` (AccessDenied), `PermissionError`, `403 Forbidden`, `os.ErrPermission` | `ACCESS_DENIED` | Security |
| `botocore.exceptions.CredentialRetrievalError`, `ExpiredSignatureError`, `jwt.exceptions.InvalidTokenError`, `jose.exceptions.*` | `AUTH_EXPIRED` | Security |
| `watchtower.CloudWatchLogBatchError`, `boto3.exceptions.*`, `google.api_core.exceptions.*` | `CLOUD_SERVICE_ERROR` | Infrastructure |
| `redis.exceptions.ConnectionError`, `redis.exceptions.TimeoutError`, `redis.exceptions.BusyLoadingError` | `CACHE_ERROR` | Cache |
| `rq.exceptions.DequeuingError`, `celery.exceptions.TimeoutError`, `pika.exceptions.AMQPConnectionError`, `github.com/hibiken/asynq.*` | `TASK_QUEUE_ERROR` | Infrastructure |
| `httpx.HTTPStatusError` (429), `ratelimit.exception.*` | `RATE_LIMIT_EXCEEDED` | Traffic |
| `httpx.HTTPStatusError` (5xx), `requests.exceptions.HTTPError` (5xx), `pywa.errors.WhatsAppError` | `UPSTREAM_FAILURE` | Integration |
| `httpx.TimeoutException`, `requests.exceptions.Timeout`, `context.DeadlineExceeded` | `INTEGRATION_TIMEOUT` | Upstream |
| `net.OpError`, `net.AddrError`, `syscall.ECONNREFUSED`, `syscall.ETIMEDOUT` | `NETWORK_ERROR` | System/OS |
| `json.JSONDecodeError`, `yaml.YAMLError`, `marshmallow.*`, `weasyprint.WeasyPrintException`, `pandas.errors.*`, `numpy.linalg.LinAlgError` | `DATA_MALFORMED` | Logic |
| `ZeroDivisionError`, `KeyError`, `AttributeError`, `RecursionError`, `Panic` | `LOGIC_CRASH` | System |
| `FileNotFoundError`, `IOError`, `os.ErrClosed` | `F_SYSTEM_ERROR` | System |

### 🏷️ 3.2 Dynamic Categorization (How Developers Add New Ones)
- **Method A (Inheritance)**:
  ```python
  class InventorySyncError(HeimdallBaseError):
      category = "INVENTORY_SYNC_FAILED"
  ```
- **Method B (Context Manager)**:
  ```python
  with heimdall.hint("GOWAY_API_FAIL"):
      call_goway_api()
  ```

---

## 🏗️ Epic 4: The Bifrost Local (Local Orchestration & DX)
**Goal**: Seamless developer workflow and manual verification on local machines.

### 💻 4.1 Local Shared Library Linking
- **Local Build**: Compile `libheimdall.so` into a known local directory (e.g., `./bin/`).
- **Env Configuration**: Document how to set `LD_LIBRARY_PATH` (Linux/Mac) or `PATH` (Windows) locally so the Python/Go adapters can find the library without specific Docker setup.

### 🛠️ 4.2 Manual Pilot Integration
- **Manual Adapter Install**: Install the local `zfw-heimdall-python` via `pip install -e` (editable mode).
- **Middleware Hook**: Manually add the middleware to `zorms/main.py`.
- **The "Kill Switch" Test**: Prove that Heimdall won't take down the service if things go wrong (e.g., renaming the `.so` file while running).

---

## 🚀 Epic 5: The Nine Realms Rollout (Service Integration)
**Goal**: Activating Heimdall in the wild after successful local pilot.

### 🧪 5.1 Pilot: `zorms`
1.  Identify the unique GoWay/EasyEcom integrations in `zorms` and wrap them in the `with heimdall.hint(...)` context manager.
2.  Run the service locally, trigger high-frequency errors (404s, 500s), and verify the JSON output matches the "Premium" spec in the terminal.
3.  Remove the existing `try/except` in `RequestLogMiddleware` that returns `str(e)`—Heimdall now owns that.

### 📊 5.2 Dashboard: `zfw_dashboard`
1.  Update the Axios/Fetch interceptor to correctly parse the new "Premium JSON" format and show the user the "Reference ID".

### 🔗 5.3 Shared Service: `z-or-cs`
1.  Apply the `zfw-heimdall-go` adapter.
2.  Verify that `trace_id` headers are passed when `zorms` calls `z-or-cs` in the local dev environment.

---

## 🛡️ Epic 6: The Ragnarök Defense (Failure & Stress Testing)
**Goal**: Proving the system is bulletproof.

- **Scenario A (Logic Crash)**: Force a `ZeroDivisionError` in a core route. Verify the user sees JSON with a `LOGIC_CRASH` category + Trace ID, NO stack trace.
- **Scenario B (Header Leak)**: Verify that headers like `Authorization` are untouched by the middleware.
- **Scenario C (Handler Death)**: Intentionally delete `libheimdall.so` from the container. Verify the Python Adapter "Fails-Open" and still returns a safe JSON error instead of crashing the process.

---

## 🛑 What NOT to build
- **Do NOT** wrap internal health checks (`/health`, `/ready`).
- **Do NOT** wrap internal metrics endpoints (`/metrics`).
- **Do NOT** attempt to "fix" the error logs; Heimdall only sanitizes the *Response* and the *Sanitized Log*; the messy original logs stay same for deep debugging.

