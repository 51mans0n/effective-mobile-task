// Package middleware defines custom HTTP middleware for logging, recovery, etc.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// ZapLogger is a chi middleware that logs requests using zap.
func ZapLogger(l *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			reqID := middleware.GetReqID(r.Context())

			next.ServeHTTP(ww, r)

			l.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.Status()),
				zap.Int("bytes", ww.BytesWritten()),
				zap.Int64("latency_ms", time.Since(start).Milliseconds()),
				zap.String("request_id", reqID),
				zap.String("remote", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)
		})
	}
}
