package slog

import (
	"log/slog"
	"runtime"
	"testing"
)

func TestLogger_Slog_1(t *testing.T) {
	l := New(slog.Default(), slog.LevelInfo,
		slog.String("os", runtime.GOOS),
		slog.String("arch", runtime.GOARCH),
	)
	l.Printf("Hello, %s", "World")
}
