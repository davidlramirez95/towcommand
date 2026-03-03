package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

// newTestLogger creates a logger that writes JSON to a buffer for inspection.
func newTestLogger(level string) (*slog.Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	lvl := parseLevel(level)
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: lvl})
	handler2 := handler.WithAttrs([]slog.Attr{
		slog.String("stage", "dev"),
		slog.String("function_name", "test-func"),
		slog.String("version", "1"),
	})
	return slog.New(&contextHandler{inner: handler2}), buf
}

func parseLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse log line: %v\nraw: %s", err, buf.String())
	}
	return m
}

func TestNew_OutputsValidJSON(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")
	logger.Info("test message", "key", "value")

	m := parseLogLine(t, buf)
	if m["msg"] != "test message" {
		t.Errorf("msg = %v, want %q", m["msg"], "test message")
	}
	if m["key"] != "value" {
		t.Errorf("key = %v, want %q", m["key"], "value")
	}
}

func TestNew_IncludesDefaultAttributes(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")
	logger.Info("test")

	m := parseLogLine(t, buf)
	if m["stage"] != "dev" {
		t.Errorf("stage = %v, want %q", m["stage"], "dev")
	}
	if m["function_name"] != "test-func" {
		t.Errorf("function_name = %v, want %q", m["function_name"], "test-func")
	}
	if m["version"] != "1" {
		t.Errorf("version = %v, want %q", m["version"], "1")
	}
}

func TestNew_RespectsLogLevel(t *testing.T) {
	tests := []struct {
		name         string
		level        string
		logFn        func(*slog.Logger)
		expectOutput bool
	}{
		{"INFO filters DEBUG", "INFO", func(l *slog.Logger) { l.Debug("debug msg") }, false},
		{"INFO allows INFO", "INFO", func(l *slog.Logger) { l.Info("info msg") }, true},
		{"WARN filters INFO", "WARN", func(l *slog.Logger) { l.Info("info msg") }, false},
		{"WARN allows WARN", "WARN", func(l *slog.Logger) { l.Warn("warn msg") }, true},
		{"ERROR filters WARN", "ERROR", func(l *slog.Logger) { l.Warn("warn msg") }, false},
		{"ERROR allows ERROR", "ERROR", func(l *slog.Logger) { l.Error("err msg") }, true},
		{"DEBUG allows DEBUG", "DEBUG", func(l *slog.Logger) { l.Debug("debug msg") }, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := newTestLogger(tt.level)
			tt.logFn(logger)
			hasOutput := buf.Len() > 0
			if hasOutput != tt.expectOutput {
				t.Errorf("output = %v, want %v (buf: %s)", hasOutput, tt.expectOutput, buf.String())
			}
		})
	}
}

func TestContextHandler_ExtractsCorrelationID(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	ctx := context.Background()
	ctx = SetCorrelationID(ctx, "corr-123")

	logger.InfoContext(ctx, "with correlation")

	m := parseLogLine(t, buf)
	if m["correlation_id"] != "corr-123" {
		t.Errorf("correlation_id = %v, want %q", m["correlation_id"], "corr-123")
	}
}

func TestContextHandler_ExtractsBookingID(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	ctx := context.Background()
	ctx = SetBookingID(ctx, "booking-456")

	logger.InfoContext(ctx, "with booking")

	m := parseLogLine(t, buf)
	if m["booking_id"] != "booking-456" {
		t.Errorf("booking_id = %v, want %q", m["booking_id"], "booking-456")
	}
}

func TestContextHandler_ExtractsUserID(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	ctx := context.Background()
	ctx = SetUserID(ctx, "user-789")

	logger.InfoContext(ctx, "with user")

	m := parseLogLine(t, buf)
	if m["user_id"] != "user-789" {
		t.Errorf("user_id = %v, want %q", m["user_id"], "user-789")
	}
}

func TestContextHandler_AllContextValues(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	ctx := context.Background()
	ctx = SetCorrelationID(ctx, "corr-1")
	ctx = SetBookingID(ctx, "book-2")
	ctx = SetUserID(ctx, "user-3")

	logger.InfoContext(ctx, "full context")

	m := parseLogLine(t, buf)
	if m["correlation_id"] != "corr-1" {
		t.Errorf("correlation_id = %v, want %q", m["correlation_id"], "corr-1")
	}
	if m["booking_id"] != "book-2" {
		t.Errorf("booking_id = %v, want %q", m["booking_id"], "book-2")
	}
	if m["user_id"] != "user-3" {
		t.Errorf("user_id = %v, want %q", m["user_id"], "user-3")
	}
}

func TestContextHandler_NoContextValues(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	logger.InfoContext(context.Background(), "no context")

	m := parseLogLine(t, buf)
	if _, ok := m["correlation_id"]; ok {
		t.Error("correlation_id should not be present when not set in context")
	}
	if _, ok := m["booking_id"]; ok {
		t.Error("booking_id should not be present when not set in context")
	}
	if _, ok := m["user_id"]; ok {
		t.Error("user_id should not be present when not set in context")
	}
}

func TestWithContext_AddsContextAttrs(t *testing.T) {
	logger, buf := newTestLogger("DEBUG")

	ctx := context.Background()
	ctx = SetCorrelationID(ctx, "corr-wc")
	ctx = SetUserID(ctx, "user-wc")

	ctxLogger := WithContext(ctx, logger)
	ctxLogger.Info("via WithContext")

	m := parseLogLine(t, buf)
	if m["correlation_id"] != "corr-wc" {
		t.Errorf("correlation_id = %v, want %q", m["correlation_id"], "corr-wc")
	}
	if m["user_id"] != "user-wc" {
		t.Errorf("user_id = %v, want %q", m["user_id"], "user-wc")
	}
}

func TestWithContext_EmptyContext(t *testing.T) {
	logger, _ := newTestLogger("DEBUG")

	ctxLogger := WithContext(context.Background(), logger)
	if ctxLogger != logger {
		t.Error("WithContext should return same logger when context has no values")
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"DEBUG", slog.LevelDebug},
		{"debug", slog.LevelDebug},
		{"INFO", slog.LevelInfo},
		{"info", slog.LevelInfo},
		{"WARN", slog.LevelWarn},
		{"WARNING", slog.LevelWarn},
		{"ERROR", slog.LevelError},
		{"unknown", slog.LevelInfo},
		{"", slog.LevelInfo},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := parseLevel(tt.input); got != tt.want {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNew_CreatesWorkingLogger(t *testing.T) {
	// Just verify New() doesn't panic and returns non-nil.
	logger := New("dev", "my-func", "v1", "DEBUG")
	if logger == nil {
		t.Fatal("New() returned nil")
	}
}
