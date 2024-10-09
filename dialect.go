package hyperbun

import (
	"github.com/blink-io/hypersql"
)

const (
	DialectPostgres = hypersql.DialectPostgres

	DialectMySQL = hypersql.DialectMySQL

	DialectSQLite = hypersql.DialectSQLite
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
