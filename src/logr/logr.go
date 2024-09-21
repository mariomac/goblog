package logr

import (
	"log/slog"
	"os"
)

const (
	FieldComponent = "component"
)

func Init(level slog.Level) {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: level <= slog.LevelDebug,
		Level:     level,
	})))
}

func Get() *slog.Logger {
	return slog.Default()
}
