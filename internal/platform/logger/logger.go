// Package logger provides structured JSON logging using slog.
// It follows 12-Factor XI: logs as event streams to stdout → CloudWatch.
package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type contextKey string

const (
	// CorrelationIDKey is the context key for correlation IDs.
	CorrelationIDKey contextKey = "correlation_id"
	BookingIDKey     contextKey = "booking_id"
	UserIDKey        contextKey = "user_id"
)

// New creates a configured slog.Logger with JSON output to stdout.
// Default attributes (stage, function_name, version) are baked into every log line.
// Log level is controlled by the logLevel parameter (DEBUG, INFO, WARN, ERROR).
func New(stage, functionName, version, logLevel string) *slog.Logger {
	level := parseLevel(logLevel)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	// Wrap with default attributes
	handler2 := handler.WithAttrs([]slog.Attr{
		slog.String("stage", stage),
		slog.String("function_name", functionName),
		slog.String("version", version),
	})

	return slog.New(&contextHandler{inner: handler2})
}

// WithContext returns a logger that extracts correlation_id, booking_id, and user_id
// from the context and includes them in every log line.
func WithContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	attrs := extractContextAttrs(ctx)
	if len(attrs) == 0 {
		return logger
	}
	return logger.With(attrsToAny(attrs)...)
}

// SetCorrelationID returns a context with the given correlation ID.
func SetCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}

// SetBookingID returns a context with the given booking ID.
func SetBookingID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, BookingIDKey, id)
}

// SetUserID returns a context with the given user ID.
func SetUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, UserIDKey, id)
}

func parseLevel(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func extractContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	if v, ok := ctx.Value(CorrelationIDKey).(string); ok && v != "" {
		attrs = append(attrs, slog.String("correlation_id", v))
	}
	if v, ok := ctx.Value(BookingIDKey).(string); ok && v != "" {
		attrs = append(attrs, slog.String("booking_id", v))
	}
	if v, ok := ctx.Value(UserIDKey).(string); ok && v != "" {
		attrs = append(attrs, slog.String("user_id", v))
	}
	return attrs
}

func attrsToAny(attrs []slog.Attr) []any {
	result := make([]any, len(attrs))
	for i, a := range attrs {
		result[i] = a
	}
	return result
}

// contextHandler is a slog.Handler that automatically extracts context values.
type contextHandler struct {
	inner slog.Handler
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

//nolint:gocritic // slog.Handler interface requires value receiver
func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := extractContextAttrs(ctx)
	for _, a := range attrs {
		r.AddAttrs(a)
	}
	return h.inner.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{inner: h.inner.WithGroup(name)}
}
