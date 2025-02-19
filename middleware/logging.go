package middleware

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
	"github.com/google/uuid"
)

type LoggingMiddleware struct {
	logger zerolog.Logger
	next   http.Handler
}

func NewLoggingMiddleware(next http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: log.Logger,
		next:   next,
	}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestID := uuid.New().String()
	w.Header().Set("X-Request-ID", requestID)

	wrapped := wrapResponseWriter(w)

	defer func() {
		duration := time.Since(start)

		m.logger.Info().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_ip", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Int("status_code", wrapped.status).
			Dur("duration", duration).
			Msg("request completed")

		if err := recover(); err != nil {
			m.logger.Error().
				Str("request_id", requestID).
				Interface("error", err).
				Msg("request panic")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()

	m.next.ServeHTTP(wrapped, r)
}

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
	if !rw.wroteHeader {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
		rw.wroteHeader = true
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}