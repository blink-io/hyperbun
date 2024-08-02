//go:build !nosqlite

package hyperbun

import (
	"context"

	"github.com/blink-io/hyperbun/dialect/sqlitedialect"
)

func NewSQLiteDialect(ctx context.Context, ops ...DialectOption) Dialect {
	dopts := applyDialectOptions(ops...)
	sops := make([]sqlitedialect.Option, 0)
	if dopts.loc != nil {
		sops = append(sops, sqlitedialect.Location(dopts.loc))
	}
	return sqlitedialect.New(sops...)
}
