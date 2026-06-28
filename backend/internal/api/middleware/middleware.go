// Package middleware provides HTTP middleware for the API: structured request
// logging, panic recovery, and injection of a request-scoped logger into the
// context so handlers and the service layer can log with correlation IDs.
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/JesusCabreraReveles/h3-explorer-ui/backend/pkg/logging"
)

// RequestLogger logs one structured line per request and seeds the context with
// a logger pre-tagged with the request ID, method, and path. It relies on
// chi's RequestID middleware having run first to populate the request ID.
func RequestLogger(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := middleware.GetReqID(r.Context())
			logger := base.With(
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			// Make the tagged logger available to downstream layers.
			ctx := logging.IntoContext(r.Context(), logger)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r.WithContext(ctx))

			logger.Info("request completed",
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

// Recoverer converts panics into 500 responses and logs them with a stack
// trace via the request-scoped logger, preventing a single bad request from
// taking down the server.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// http.ErrAbortHandler is a sentinel the stdlib uses to abort a
				// handler intentionally; propagate it untouched. recover() yields
				// an `any`, so the identity comparison is the documented pattern
				// (mirrors net/http's own recovery).
				if rec == http.ErrAbortHandler { //nolint:errorlint // rec is any, sentinel identity check
					panic(rec)
				}
				logging.FromContext(r.Context()).Error("recovered from panic",
					slog.Any("panic", rec),
				)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":{"code":"internal_error","message":"internal server error"}}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
