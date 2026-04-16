package heimdall

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMissingTraceIDGeneratesNew(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			t.Error("expected trace id but got empty")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := HeimdallHandler(next)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestExistingTraceIDIsPreserved(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-ID")
		if traceID != "existing-id" {
			t.Errorf("expected existing-id, got %s", traceID)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := HeimdallHandler(next)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Trace-ID", "existing-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
}
