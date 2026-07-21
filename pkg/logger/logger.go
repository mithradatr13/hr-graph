package logger

import (
	"context"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		_ = h.Handle(ctx, r)
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return &MultiHandler{handlers: handlers}
}

func InitLogger() *slog.Logger {
	// General Log File with Rotation
	generalLogRotator := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100, // Megabytes
		MaxBackups: 3,
		MaxAge:     28,   // Days
		Compress:   true, // Disabled by default, true for gzip
	}

	// Error Log File with Rotation
	errorLogRotator := &lumberjack.Logger{
		Filename:   "logs/error.log",
		MaxSize:    50,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	// Console writer (Stdout)
	consoleWriter := os.Stdout

	// Filter for errors only
	errorLevelFilter := func(groups []string, a slog.Attr) slog.Attr {
		return a
	}

	// Create JSON handlers for files and text/json for console
	jsonOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	consoleHandler := slog.NewJSONHandler(consoleWriter, jsonOpts)
	fileHandler := slog.NewJSONHandler(generalLogRotator, jsonOpts)

	// Error handler logs only ERROR and above to error.log
	errorOpts := &slog.HandlerOptions{
		Level:       slog.LevelError,
		ReplaceAttr: errorLevelFilter,
	}
	errFileHandler := slog.NewJSONHandler(errorLogRotator, errorOpts)

	// Combine handlers using MultiHandler
	multiHandler := &MultiHandler{
		handlers: []slog.Handler{
			consoleHandler,
			fileHandler,
			errFileHandler,
		},
	}

	logger := slog.New(multiHandler)
	slog.SetDefault(logger)

	return logger
}
