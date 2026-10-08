package httpapi

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder wraps http.ResponseWriter to remember the status code,
type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader saves the status code and then writes it as usual.
func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// logRequests logs one line per request with its method, path, status
// and duration.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

// recoverPanic turns a panic in a handler into a JSON 500 response
// instead of a dropped connection.
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", "err", err)
				writeJSON(w, http.StatusInternalServerError, ErrorResponse{
					Error: ErrorBody{Code: "INTERNAL_ERROR", Message: "internal server error"},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}