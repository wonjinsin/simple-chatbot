package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/wonjinsin/go-boilerplate/pkg/constants"
)

// HTTPLogger logs HTTP requests with TrID
func HTTPLogger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap ResponseWriter to capture status code and bytes
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Get TrID from context
			trID := ""
			if ctx := r.Context(); ctx != nil {
				if id, ok := ctx.Value(constants.ContextKeyTrID).(string); ok {
					trID = id
				}
			}

			// Process request
			next.ServeHTTP(ww, r)

			// Log after request is processed
			duration := time.Since(start)

			log.Info().
				Str("trid", trID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration_ms", duration).
				Msg("http request")
		})
	}
}
