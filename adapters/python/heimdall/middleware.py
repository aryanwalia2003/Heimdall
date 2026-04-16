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
    # 1. Try to find the library inside the package folder (Self-Bundled)
    current_dir = os.path.dirname(__file__)
    internal_lib_path = os.path.join(current_dir, "libheimdall.so")
    
    # 2. Fall back to environment variable if not found internally
    lib_path = os.environ.get("HEIMDALL_LIB_PATH")
    
    if os.path.exists(internal_lib_path):
        _lib = ctypes.CDLL(internal_lib_path)
    elif lib_path:
        _lib = ctypes.CDLL(lib_path)
    else:
        # Last resort: try default name in bin/
        _lib = ctypes.CDLL("bin/libheimdall.so")
    _lib.ProcessError.argtypes = [ctypes.c_char_p]
    _lib.ProcessError.restype = ctypes.c_void_p
    _lib.FreeString.argtypes = [ctypes.c_void_p]
except Exception as e:
    print(f"Heimdall warning: Could not load shared library: {e}", file=sys.stderr)

import logging

logger = logging.getLogger("heimdall")
logger.setLevel(logging.ERROR)
if not logger.handlers:
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(logging.Formatter('[%(name)s] %(levelname)s - %(message)s'))
    logger.addHandler(handler)

def _process_with_heimdall(exc: Exception) -> JSONResponse:
    trace = traceback.format_exc()
    category = getattr(exc, "heimdall_category", getattr(exc, "category", None))
    trace_id = current_trace_id.get()
    
    # Log the full technical truth on the server for correlation
    logger.error(
        f"TraceID={trace_id} | Exception={type(exc).__name__}: {str(exc)}\n"
        f"--- Full Stack Trace ---\n{trace}\n------------------------"
    )

    payload = {
        "trace_id": trace_id,
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
