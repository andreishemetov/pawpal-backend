package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/andreishemetov/pawpal/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(lrw, r)

		status := strconv.Itoa(lrw.status)

		metrics.HTTPRequestsTotal.
			WithLabelValues(r.Method, r.URL.Path, status).
			Inc()

		metrics.HTTPRequestDuration.
			WithLabelValues(r.Method, r.URL.Path).
			Observe(time.Since(start).Seconds())
	})
}
