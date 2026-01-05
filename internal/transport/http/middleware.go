package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type contextKey string

const requestIDKey contextKey = "request_id"

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"route", "method", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route", "method"},
	)

	aiRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ai_requests_total",
			Help: "Total number of AI requests",
		},
		[]string{"status"},
	)

	aiRequestDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "ai_request_duration_seconds",
			Help:    "AI request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Context().Value(requestIDKey)
		if requestID == nil {
			requestID = "unknown"
		}

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		logEntry := map[string]interface{}{
			"level":      "info",
			"msg":        "HTTP request",
			"request_id": requestID,
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     ww.Status(),
			"duration_ms": duration.Milliseconds(),
		}
		if logJSON, err := json.Marshal(logEntry); err == nil {
			log.Println(string(logJSON))
		} else {
			log.Printf("HTTP request: request_id=%v method=%s path=%s status=%d duration_ms=%d",
				requestID, r.Method, r.URL.Path, ww.Status(), duration.Milliseconds())
		}
	})
}

// MetricsMiddleware records Prometheus metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		route := r.URL.Path
		method := r.Method

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		status := http.StatusText(ww.Status())

		httpRequestsTotal.WithLabelValues(route, method, status).Inc()
		httpRequestDuration.WithLabelValues(route, method).Observe(duration.Seconds())
	})
}

// RecordAIRequest records metrics for an AI request
func RecordAIRequest(duration time.Duration, status string) {
	aiRequestsTotal.WithLabelValues(status).Inc()
	aiRequestDuration.Observe(duration.Seconds())
}

