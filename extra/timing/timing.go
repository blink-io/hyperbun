package timing

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

type (
	hook struct {
		logf func(string, ...any)
	}
)

var _ bun.QueryHook = (*hook)(nil)

func New(ops ...Option) bun.QueryHook {
	h := new(hook)
	for _, o := range ops {
		o(h)
	}
	if h.logf == nil {
		h.logf = func(format string, args ...any) {
			slog.Info(fmt.Sprintf(format, args...))
		}
	}
	return h
}

func (h *hook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *hook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	h.logf("[QueryHook][Timing] Executed SQL, timing cost [%s] for: %s", time.Since(event.StartTime), event.Query)
}
