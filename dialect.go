package hyperbun

import (
	"github.com/blink-io/hypersql"
)

func GetDialect(dialect string, ops ...DialectOption) (Dialect, error) {
	switch hypersql.GetFormalDialect(dialect) {
	case hypersql.DialectPostgres:
		return NewPostgresDialect(ops...), nil
	case hypersql.DialectMySQL:
		return NewMySQLDialect(ops...), nil
	case hypersql.DialectSQLite:
		return NewSQLiteDialect(ops...), nil
	default:
		return nil, hypersql.ErrUnsupportedDialect
	}
}
