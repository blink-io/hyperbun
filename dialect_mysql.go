package hyperbun

import (
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func NewMySQLDialect(ops ...DialectOption) Dialect {
	return mysqldialect.New()
}
