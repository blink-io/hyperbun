package hyperbun

import (
	"context"

	"github.com/uptrace/bun/dialect/mysqldialect"
)

func NewMySQLDialect(ctx context.Context, ops ...DialectOption) Dialect {
	return mysqldialect.New()
}
