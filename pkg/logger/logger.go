package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var (
	rootLogger      *slog.Logger
	requestIDHeader string
	fileOut         *os.File
)

// Init initialises the package logger with the provided Config.
func Init(cfg Config) error {
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.Output == "" {
		cfg.Output = "stdout"
	}
	requestIDHeader = cfg.RequestIDHeader

	var w *os.File = os.Stdout
	if strings.ToLower(cfg.Output) != "stdout" {
		f, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("open log file: %w", err)
		}
		fileOut = f
		w = f
	}

	suppressKeys := func(groups []string, a slog.Attr) slog.Attr {

		for _, g := range groups {
			if g == "client" {
				return slog.Attr{}
			}
		}
		return a
	}

	opts := &slog.HandlerOptions{}

	switch strings.ToLower(cfg.Level) {
	case "debug":
		opts.Level = slog.LevelDebug
		opts = &slog.HandlerOptions{AddSource: true}
	case "info":
		opts = &slog.HandlerOptions{AddSource: false, ReplaceAttr: suppressKeys}
		opts.Level = slog.LevelInfo
	case "warn", "warning":
		opts = &slog.HandlerOptions{AddSource: false, ReplaceAttr: suppressKeys}
		opts.Level = slog.LevelWarn
	case "error":
		opts = &slog.HandlerOptions{AddSource: false, ReplaceAttr: suppressKeys}
		opts.Level = slog.LevelError
	default:
		opts = &slog.HandlerOptions{AddSource: false, ReplaceAttr: suppressKeys}
		opts.Level = slog.LevelInfo
	}

	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "text", "standard", "default", "console":
		handler = slog.NewTextHandler(w, opts)
	case "json":
		handler = slog.NewJSONHandler(w, opts)
	default:
		// fall back to text for unknown values to keep logs readable
		handler = slog.NewTextHandler(w, opts)
	}

	rootLogger = slog.New(handler)
	slog.SetDefault(rootLogger)
	return nil
}

// InitFromEnv reads environment variables and inits the logger.
func InitFromEnv() error {
	cfg := DefaultConfig()
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.Level = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.Format = v
	}
	if v := os.Getenv("LOG_OUTPUT"); v != "" {
		cfg.Output = v
	}
	if v := os.Getenv("REQUEST_ID_HEADER"); v != "" {
		cfg.RequestIDHeader = v
	}
	return Init(cfg)
}

// Sync flushes/closes any resources used by the logger (e.g. file handles).
func Sync() error {
	if fileOut != nil {
		err := fileOut.Close()
		fileOut = nil
		return err
	}
	return nil
}

// WithAttrs returns a child logger with additional attributes. It accepts any
// arguments supported by slog.Logger.With (convenience wrapper).
func WithAttrs(attrs ...any) *slog.Logger {
	return slog.Default().With(attrs...)
}

type ctxKey struct{}

// NewContext returns a context that carries the provided logger.
func NewContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the logger from the context if present, otherwise the default logger.
func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.Default()
	}
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(*slog.Logger); ok {
			return l
		}
	}
	return slog.Default()
}
