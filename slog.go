package midas

import (
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
)

// SetupSlog Configure the default [slog.Logger] with midas defaults. Available variables:
//   - MIDAS_LOG_LEVEL: control slog's default log level, one of [debug, info, warn, error], default to info
//   - MIDAS_LOG_ADD_SOURCE: control if the source must be logged, default to false
func SetupSlog(cfg Config) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: cfg.GetBool("MIDAS_LOG_ADD_SOURCE", false),
		Level:     parseSlogLevel(cfg.Get("MIDAS_LOG_LEVEL", "info")),
	}))
	logger = logger.With("executionId", uuid.NewString())

	slog.SetDefault(logger)
}

func parseSlogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	return slog.LevelInfo
}
