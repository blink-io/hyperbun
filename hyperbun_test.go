package hyperbun

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestSafeQuery_1(t *testing.T) {
	q := "abc"
	ss := doSafeQuery(q, "a", "b", "c")
	require.NotNil(t, ss)
}

func TestBun_Logger_1(t *testing.T) {
	bun.SetLogger(Logf(func(format string, args ...any) {
		msg := fmt.Sprintf(format, args...)
		slog.Info(msg)
	}))
}

func TestTime_Format_1(t *testing.T) {
	tt := time.Now()
	tstr := tt.Format(time.RFC3339Nano)

	fmt.Printf("%s\n", tstr)
}
