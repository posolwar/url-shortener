package config

import (
	"io"
	"log/slog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string, w io.Writer) {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envDev:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug})
	case envProd:
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	slog.SetDefault(slog.New(handler))
}
