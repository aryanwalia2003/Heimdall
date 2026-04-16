import traceback
import contextvars
import ulid
from fastapi import Request
from fastapi.responses import JSONResponse
from starlette.middleware.base import BaseHTTPMiddleware

current_trace_id = contextvars.ContextVar("trace_id", default=None)

def _get_or_create_trace_id(request: Request) -> str:
    trace_id = request.headers.get("X-Trace-ID")
    if not trace_id:
        trace_id = str(ulid.new())
    current_trace_id.set(trace_id)
    return trace_id

def _process_with_heimdall(exc: Exception) -> JSONResponse:
    trace = traceback.format_exc()
    try:
        # In a real deployed setup we'd call libheimdall.so here.
        # Since we're missing it, it raises an exception which fail-opens.
        raise RuntimeError("libheimdall missing")
    except Exception:
        return JSONResponse(
            status_code=500,
            content={
                "success": False,
                "message": "Critical System Error",
                "trace_id": current_trace_id.get()
            }
        )

class HeimdallMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        _get_or_create_trace_id(request)
        try:
            return await call_next(request)
        except Exception as e:
            return _process_with_heimdall(e)
