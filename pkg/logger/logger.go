package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
)

func InitLogger(cfg *configs.AppConfig) {
	level := getLogLevel(cfg.LogLevel)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Env == "dev" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(
		slog.String("service", cfg.Name),
		slog.String("env", cfg.Env),
	)

	slog.SetDefault(logger)
}

func getLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
