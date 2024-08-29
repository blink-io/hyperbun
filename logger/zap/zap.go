package zap

import (
	"fmt"

	"github.com/uptrace/bun"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Printf(format string, args ...interface{})
}

type logger struct {
	zl *zap.Logger
	lv zapcore.Level
	fd []zapcore.Field
}

func New(zl *zap.Logger, lv zapcore.Level, fd ...zap.Field) Logger {
	return &logger{
		zl: zl,
		lv: lv,
		fd: fd,
	}
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.zl.Log(l.lv, fmt.Sprintf(format, args...), l.fd...)
}

func SetLogger(zl *zap.Logger, lv zapcore.Level, fd ...zap.Field) {
	bun.SetLogger(New(zl, lv, fd...))
}
