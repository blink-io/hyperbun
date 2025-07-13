package slog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun"
)

type Logger interface {
	Printf(format string, args ...any)
}

type logger struct {
	sl *slog.Logger
	lv slog.Leveler
	as []slog.Attr
}

func New(sl *slog.Logger, lv slog.Leveler, as ...slog.Attr) Logger {
	return &logger{
		sl: sl,
		lv: lv,
		as: as,
	}
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.sl.LogAttrs(context.Background(), l.lv.Level(), fmt.Sprintf(format, args...), l.as...)
}

func SetLogger(sl *slog.Logger, lv slog.Leveler, as ...slog.Attr) {
	bun.SetLogger(New(sl, lv, as...))
}
