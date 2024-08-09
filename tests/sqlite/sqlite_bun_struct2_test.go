package sqlite

import (
	"testing"

	bunx "github.com/blink-io/hyperbun"

	"github.com/stretchr/testify/require"
)

func TestSqlite_Bun_Tuple2SQL_1(t *testing.T) {
	db := getSqliteDB()

	vals, err := bunx.TypeTuple2SQL[int64, string](ctx, db, "select id, guid from users limit ?", 10)
	require.NoError(t, err)
	require.NotNil(t, vals)
}

func TestSqlite_Bun_Tuple3SQL_1(t *testing.T) {
	db := getSqliteDB()

	vals, err := bunx.TypeTuple3SQL[int64, string, string](ctx, db, "select id, guid, profile from users limit ?", 10)
	require.NoError(t, err)
	require.NotNil(t, vals)
}
