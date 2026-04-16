import uvicorn
from fastapi import FastAPI, HTTPException
from heimdall.middleware import HeimdallMiddleware

app = FastAPI(title="Heimdall Pilot Service")

# Add Heimdall Middleware
app.add_middleware(HeimdallMiddleware)

@app.get("/")
async def root():
    return {"message": "Heimdall Pilot is running"}

@app.get("/trigger-error")
async def trigger_error():
    # Trigger a sample error that Heimdall should catch and sanitize
    raise ValueError("Sensitive data leak: password=secret123 and database_url=mysql://user:pass@localhost:3306/db")

@app.get("/devious")
async def devious():
    # Sensitive connection string and API key
    raise ValueError("Critical Failure: Connection to postgresql://admin:P@ssword123@zorms-db:5432/dev failed with API_KEY=AKIA_MOCK_123456789")

@app.get("/div0")
async def div0():
    return 1 / 0

@app.get("/attr")
async def attr():
    return None.some_method()

@app.get("/trigger-runtime-error")
async def trigger_runtime_error():
    raise RuntimeError("A generic runtime error occurred")

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
