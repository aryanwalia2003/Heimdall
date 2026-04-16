package logger

import (
	"log/slog"
	"os"
)

// Logger is a wrapper around slog.Logger
type Logger struct {
	*slog.Logger
}

// New creates a new structured logger
func New() *Logger {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	return &Logger{slog.New(handler)}
}

// WithTrace adds a trace_id to the logger context
func (l *Logger) WithTrace(traceID string) *Logger {
	return &Logger{l.With("trace_id", traceID)}
}
