package hyperbun

import (
	"context"

	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewPostgresDialect(ctx context.Context, ops ...DialectOption) Dialect {
	return pgdialect.New()
}
