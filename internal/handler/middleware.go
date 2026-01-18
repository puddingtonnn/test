package handler

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.PathValue("id"),
			"uri", r.URL.RequestURI(),
			"status", lrw.statusCode,
			"duration", duration.String(),
			"remote_addr", r.RemoteAddr,
		)
	})
}
