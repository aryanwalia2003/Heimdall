package heimdall

import (
	"fmt"
	"net/http"

	"github.com/aryanwalia/heimdall/pkg/wrap"
	"github.com/google/uuid"
)

// HeimdallHandler wraps an http.Handler with trace injection and panic recovery.
func HeimdallHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := getOrCreateTraceID(r)
		defer recoverAndProcessError(w, traceID)
		next.ServeHTTP(w, r)
	})
}

func getOrCreateTraceID(r *http.Request) string {
	traceID := r.Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.New().String()
		r.Header.Set("X-Trace-ID", traceID)
	}
	return traceID
}

func recoverAndProcessError(w http.ResponseWriter, traceID string) {
	if err := recover(); err != nil {
		rawErr := fmt.Sprintf(`{"error": "%v", "trace_id": "%s"}`, err, traceID)
		premiumJSON := wrap.ProcessError(rawErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(premiumJSON))
	}
}
