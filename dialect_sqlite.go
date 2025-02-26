package hyperbun

import (
	"github.com/blink-io/hyperbun/dialect/sqlitedialect"
)

func NewSQLiteDialect(ops ...DialectOption) Dialect {
	dopts := applyDialectOptions(ops...)
	sops := make([]sqlitedialect.Option, 0)
	if dopts.loc != nil {
		sops = append(sops, sqlitedialect.WithLocation(dopts.loc))
	}
	return sqlitedialect.New(sops...)
}
