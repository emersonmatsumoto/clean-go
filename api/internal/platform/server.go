package platform

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK} // Default 200
}

func metricsMiddleware(next http.Handler) http.Handler {
	var meter = otel.Meter("http-server")
	httpDuration, _ := meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("Duração das requisições HTTP"),
		metric.WithUnit("ms"),
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds() * 1000

		attrs := metric.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
			attribute.Int("status", rw.statusCode),
		)

		httpDuration.Record(r.Context(), duration, attrs)
	})
}

func NewServer(addr string, handler http.Handler) *http.Server {
	handler = metricsMiddleware(handler)

	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
