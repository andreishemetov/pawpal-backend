package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// responseWriter is a small wrapper to capture status code and response size.
/*
http.ResponseWriter — embedded interface. This means responseWriter automatically satisfies the http.ResponseWriter interface and delegates all methods to the real writer underneath. No need to re-implement Header() etc.
status int — captures the HTTP status code when WriteHeader() is called
size int — accumulates the number of bytes written via Write()
The pattern — "embedding to override"
*/
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// ensure we record a status if none was set (default 200)
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Logging is the HTTP middleware that logs request/response info.
// It expects request-id middleware (chi/middleware.RequestID) to run earlier
// so `middleware.GetReqID(r.Context())` returns a string id.
func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}

		// Get request ID added by chi middleware.RequestID (empty if not present)
		reqID := middleware.GetReqID(r.Context())

		// Call next handler
		next.ServeHTTP(rw, r)

		// Default status in case none written
		if rw.status == 0 {
			rw.status = http.StatusOK
		}

		// Log output (plain text). Format: time method path status size duration remote addr req_id
		duration := time.Since(start)
		log.Printf(
			"%s %s %s %d %dB %s remote=%s req_id=%s",
			r.Method,
			r.Proto,
			r.URL.RequestURI(),
			rw.status,
			rw.size,
			duration.String(),
			r.RemoteAddr,
			reqID,
		)
	}
	return http.HandlerFunc(fn)
}