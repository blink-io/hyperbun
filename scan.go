package hyperbun

import (
	"context"
	"database/sql"

	"github.com/stephenafamo/scan"
)

var _ scan.Queryer = (*queryer)(nil)

func (db *DB) ScanQueryer() scan.Queryer {
	return convert(db)
}

// A Queryer that returns the concrete type [*sql.Rows]
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// convert wraps an Queryer and makes it a Queryer
func convert(wrapped Queryer) scan.Queryer {
	return queryer{wrapped: wrapped}
}

type queryer struct {
	wrapped Queryer
}

func (q queryer) QueryContext(ctx context.Context, query string, args ...any) (scan.Rows, error) {
	return q.wrapped.QueryContext(ctx, query, args...)
}
