import pytest
import traceback
from fastapi import FastAPI, Request
from fastapi.testclient import TestClient
from fastapi.responses import JSONResponse
from heimdall.middleware import HeimdallMiddleware, current_trace_id

def create_app() -> FastAPI:
    app = FastAPI()
    app.add_middleware(HeimdallMiddleware)
    
    @app.get("/success")
    async def success():
        return {"status": "ok", "trace_id": current_trace_id.get()}
        
    @app.get("/error")
    async def error():
        raise ValueError("Something went terribly wrong")
        
    return app

@pytest.fixture
def client() -> TestClient:
    return TestClient(create_app())

def test_missing_trace_id_generates_ulid(client: TestClient):
    response = client.get("/success")
    assert response.status_code == 200
    data = response.json()
    assert "trace_id" in data
    assert bool(data["trace_id"])

def test_existing_trace_id_is_preserved(client: TestClient):
    trace_id = "01ARZ3NDEKTSV4RRFFQ69G5FAV"
    response = client.get("/success", headers={"X-Trace-ID": trace_id})
    assert response.status_code == 200
    data = response.json()
    assert data["trace_id"] == trace_id

def test_exception_returns_heimdall_json(client: TestClient):
    response = client.get("/error")
    assert response.status_code == 500
    data = response.json()
    # It should have the trace_id
    assert "trace_id" in data
    assert data["success"] is False
    assert "message" in data
