package logr

import (
	"log/slog"
	"strings"
)

func ParseLevel(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	case "DEBUG":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}