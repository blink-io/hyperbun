package hyperbun

import (
	"github.com/blink-io/hypersql"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

const (
	DialectPostgres = hypersql.DialectPostgres

	DialectMySQL = hypersql.DialectMySQL

	DialectSQLite = hypersql.DialectSQLite

	DialectSQLServer = hypersql.DialectSQLServer
)

var (
	ErrUnsupportedDialect = hypersql.ErrUnsupportedDialect
	ErrUnsupportedDriver  = hypersql.ErrUnsupportedDriver
	InvalidConfig         = hypersql.InvalidConfig
)

func GetDialect(dialect string) (Dialect, error) {
	switch hypersql.GetFormalDialect(dialect) {
	case DialectPostgres:
		return pgdialect.New(), nil
	case DialectMySQL:
		return mysqldialect.New(), nil
	case DialectSQLite:
		return sqlitedialect.New(), nil
	case DialectSQLServer:
		return mssqldialect.New(), nil
	default:
		return nil, hypersql.ErrUnsupportedDialect
	}
}
