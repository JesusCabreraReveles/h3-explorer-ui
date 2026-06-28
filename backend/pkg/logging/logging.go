// Package logging centralizes structured logging setup around log/slog and
// provides helpers to carry a request-scoped logger through context. It lives
// in pkg/ because it is generic infrastructure with no domain knowledge.
package logging

import (
	"context"
	"io"
	"log/slog"
	"strings"
)

// Format selects the slog output encoding.
type Format string

const (
	// FormatJSON emits machine-parseable structured logs (production default).
	FormatJSON Format = "json"
	// FormatText emits human-friendly logs (handy for local development).
	FormatText Format = "text"
)

// ctxKey is unexported to prevent collisions with other packages' context keys.
type ctxKey struct{}

// New builds a slog.Logger writing to w at the given level and format.
func New(w io.Writer, level string, format Format) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}

	var handler slog.Handler
	if format == FormatText {
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}
	return slog.New(handler)
}

// parseLevel maps a case-insensitive string to a slog level, defaulting to info.
func parseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// IntoContext returns a copy of ctx carrying the given logger.
func IntoContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// FromContext retrieves the request-scoped logger, falling back to the default
// logger so callers never have to nil-check.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
