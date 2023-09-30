package observ

import (
	"context"
	"os"
	"sync"

	"log/slog"
)

var Logger = func() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}

var once sync.Once

var stdlog *slog.Logger

func SetLogger() {
	stdlog = Logger()
}

func Log(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

func LogError(ctx context.Context, msg string, fields ...any) {
	log(ctx, slog.LevelError, msg, fields...)
}

func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	once.Do(SetLogger)

	stdlog.Log(ctx, level, msg, args...)
}
