package logger

import (
	"context"
	"log/slog"
	"os"

	chimdlwr "github.com/go-chi/chi/v5/middleware"
)

var (
	defaultLogger   = slog.Default()
	defaultLogLevel = slog.LevelInfo
)

func InitLogger(cfglvl string) {
	var lvl slog.Level = defaultLogLevel
	if cfglvl != "" {
		if err := lvl.UnmarshalText([]byte(cfglvl)); err != nil {
			slog.Warn("failed to parse log level; using default", "err", err)
		}
	}
	var h slog.Handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
	})
	defaultLogger = slog.New(reqIDHandler{h})
}

// For general use, logs source info.
func GetLogger() *slog.Logger {
	return defaultLogger
}

// Only for router use, doesn't log source info.
func GetRequestLogger() *slog.Logger {
	return slog.New(reqIDHandler{slog.NewJSONHandler(os.Stdout, nil)})
}

type reqIDHandler struct {
	slog.Handler
}

func (h reqIDHandler) Handle(ctx context.Context, r slog.Record) error {
	id, ok := ctx.Value(chimdlwr.RequestIDKey).(string)
	if ok {
		return h.Handler.WithAttrs([]slog.Attr{slog.String("id", id)}).Handle(ctx, r)
	}
	return h.Handler.Handle(ctx, r)
}

func (h reqIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return reqIDHandler{h.Handler.WithAttrs(attrs)}
}
