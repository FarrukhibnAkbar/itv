package middleware

import (
	"log"
	"net/http"
	"time"
)

// loggingResponseWriter wraps ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter initializes a response writer wrapper
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code before passing it to the original writer
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Logging middleware logs requests with method, status code, URL, execution time, and client IP
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		duration := time.Since(start)
		clientIP := req.RemoteAddr // Extract client IP

		log.Printf("[%s] %s %d %s %s - %v", time.Now().Format(time.RFC3339), clientIP, statusCode, req.Method, req.URL.String(), duration)
	})
}
