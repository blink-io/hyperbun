package zap

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Zap_1(t *testing.T) {
	zl, err := zap.NewDevelopment()
	require.NoError(t, err)

	l := New(zl, zapcore.InfoLevel,
		zap.String("arch", runtime.GOARCH),
		zap.String("os", runtime.GOOS),
		zap.String("version", runtime.Version()),
		zap.String("compiler", runtime.Compiler),
		zap.String("goroot", runtime.GOROOT()),
	)
	l.Printf("Hello, %s", "World")
}
