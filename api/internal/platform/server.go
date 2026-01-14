package platform

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
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

var logger = global.GetLoggerProvider().Logger("http-middleware")

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

		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetBody(log.StringValue("[http.request]"))

		record.AddAttributes(
			log.String("method", r.Method),
			log.String("path", r.URL.Path),
			log.Int64("status", int64(rw.statusCode)),
			log.String("ip", r.RemoteAddr),
			log.String("user_agent", r.UserAgent()),
			log.Float64("duration_ms", float64(duration)),
		)

		if rw.statusCode >= 500 {
			record.SetSeverity(log.SeverityError)
		} else {
			record.SetSeverity(log.SeverityInfo)
		}

		logger.Emit(r.Context(), record)
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
