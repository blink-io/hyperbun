package hyperbun

import (
	"github.com/blink-io/hypersql"
)

const (
	DialectPostgres = hypersql.DialectPostgres

	DialectMySQL = hypersql.DialectMySQL

	DialectSQLite = hypersql.DialectSQLite
)

var (
	ErrUnsupportedDialect = hypersql.ErrUnsupportedDialect
	ErrUnsupportedDriver  = hypersql.ErrUnsupportedDriver
	InvalidConfig         = hypersql.InvalidConfig
)

func GetDialect(dialect string, ops ...DialectOption) (Dialect, error) {
	switch hypersql.GetFormalDialect(dialect) {
	case DialectPostgres:
		return NewPostgresDialect(ops...), nil
	case DialectMySQL:
		return NewMySQLDialect(ops...), nil
	case DialectSQLite:
		return NewSQLiteDialect(ops...), nil
	default:
		return nil, hypersql.ErrUnsupportedDialect
	}
}
