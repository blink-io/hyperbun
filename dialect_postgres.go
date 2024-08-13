package hyperbun

import (
	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewPostgresDialect(ops ...DialectOption) Dialect {
	return pgdialect.New()
}
