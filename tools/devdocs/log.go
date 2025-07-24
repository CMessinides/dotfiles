package main

import (
	"log/slog"
	"os"
)

var logLevel = new(slog.LevelVar)

func init() {
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	slog.SetDefault(slog.New(h))
}

func SetLogLevel(level slog.Level) {
	logLevel.Set(level)
}
