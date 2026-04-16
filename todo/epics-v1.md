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
| `sqlalchemy.exc.*`, `psycopg2.*`, `pymongo.*` | `DATABASE_ERROR` | Persistence |
| `fastapi.exceptions.RequestValidationError`, `pydantic.ValidationError` | `VALIDATION_ERROR` | Schema/Request |
| `requests.exceptions.Timeout`, `httpx.ConnectTimeout` | `INTEGRATION_TIMEOUT` | Upstream |
| `requests.exceptions.HTTPError` (Status 5xx from upstream) | `UPSTREAM_FAILURE` | Integration |
| `jose.exceptions.*`, `ExpiredSignatureError` | `AUTH_EXPIRED` | Security |
| `PermissionError`, `403 Forbidden` | `ACCESS_DENIED` | Security |
| `404 Not Found` | `RESOURCE_NOT_FOUND` | URL/Path |
| `ZeroDivisionError`, `KeyError`, `AttributeError` | `LOGIC_CRASH` | System |

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

## 🏗️ Epic 4: The Citadel (Deployment & Infrastructure)
**Goal**: Seamless rollout to 500+ future services using Base Images.

### 🐳 4.1 Base Image Architecture
- **Repo**: `zfw-infra/base-images`.
- **Docker-Python**:
  - `FROM python:3.11-slim`
  - `COPY --from=heimdall-builder /libheimdall.so /usr/local/lib/`
  - `ENV LD_LIBRARY_PATH=/usr/local/lib`
- **Docker-Go**:
  - Similar to Python but optimized for Go binaries.

### 🛠️ 4.2 Automated Injections
- Patch the `deploy-prod.yml` in all repos to point to `zfw-base-python` instead of generic `python:3.11`.
- This ensures that `libheimdall.so` is "just there" in the environment.

---

## 🚀 Epic 5: The Nine Realms Rollout (Service Integration)
**Goal**: Activating Heimdall in the wild.

### 🧪 5.1 Pilot: `zorms`
1.  Add `zfw-heimdall-python` to `requirements.txt`.
2.  In `main.py`: `app.add_middleware(HeimdallMiddleware)`.
3.  Identify the unique GoWay/EasyEcom integrations in `zorms` and wrap them in the `with heimdall.hint(...)` context manager.
4.  Remove the existing `try/except` in `RequestLogMiddleware` that returns `str(e)`—Heimdall now owns that.

### 📊 5.2 Dashboard: `zfw_dashboard`
1.  The `zfw_dashboard` is a frontend-heavy repo. If it has a Node/Go backend component, apply the respective adapter.
2.  If it is strictly a static React/Vue app, update the Axios/Fetch interceptor to correctly parse the new "Premium JSON" format and show the user the "Reference ID".

### 🔗 5.3 Shared Service: `z-or-cs`
1.  Apply the `zfw-heimdall-go` adapter.
2.  Verify that `trace_id` headers are passed when `zorms` calls `z-or-cs`.

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

