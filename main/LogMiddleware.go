package main

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Printf("RECOVER %v %v", err, debug.Stack())
				}
			}()

			start := time.Now()

			wrw := wrapResponseWriter(w)

			wrw.Header().Set(HTTP_VER_HEADER, VERSION)
			next.ServeHTTP(wrw, r)

			logger.Printf("%v %03d %v %v %v",
				r.RemoteAddr,
				wrw.status,
				r.Method,
				r.RequestURI,
				time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
