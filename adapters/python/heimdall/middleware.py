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
        trace_id = str(ulid.ulid())
    current_trace_id.set(trace_id)
    return trace_id

import ctypes
import json
import os
import sys

# Load the shared library
_lib = None
try:
    # Look for libheimdall.so in the bin directory relative to the project root
    # For local development, we can also check LD_LIBRARY_PATH
    lib_path = os.environ.get("HEIMDALL_LIB_PATH", "bin/libheimdall.so")
    _lib = ctypes.CDLL(lib_path)
    _lib.ProcessError.argtypes = [ctypes.c_char_p]
    _lib.ProcessError.restype = ctypes.c_void_p
    _lib.FreeString.argtypes = [ctypes.c_void_p]
except Exception as e:
    print(f"Heimdall warning: Could not load shared library: {e}", file=sys.stderr)

def _process_with_heimdall(exc: Exception) -> JSONResponse:
    trace = traceback.format_exc()
    category = getattr(exc, "heimdall_category", getattr(exc, "category", None))
    
    payload = {
        "trace_id": current_trace_id.get(),
        "error": f"{type(exc).__name__}: {str(exc)}",
        "stack_trace": trace,
        "category": category,
        "status_code": 500
    }
    
    try:
        if _lib:
            input_json = json.dumps(payload).encode('utf-8')
            output_ptr = _lib.ProcessError(input_json)
            if output_ptr:
                # Cast the void pointer to a char pointer and get the value
                ptr = ctypes.cast(output_ptr, ctypes.c_char_p)
                output_json = ptr.value.decode('utf-8')
                result = json.loads(output_json)
                _lib.FreeString(output_ptr)
                return JSONResponse(status_code=500, content=result)
        
        raise RuntimeError("libheimdall missing or failed")
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
