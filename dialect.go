package hyperbun

import (
	"context"

	"github.com/blink-io/hypersql"
)

func GetDialect(ctx context.Context, c *Config, ops ...DialectOption) (Dialect, error) {
	switch hypersql.GetFormalDialect(c.Dialect) {
	case hypersql.DialectPostgres:
		return NewPostgresDialect(ctx, ops...), nil
	case hypersql.DialectMySQL:
		return NewMySQLDialect(ctx, ops...), nil
	case hypersql.DialectSQLite:
		return NewSQLiteDialect(ctx, ops...), nil
	default:
		return nil, hypersql.ErrUnsupportedDialect
	}

}
