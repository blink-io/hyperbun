package sqlite

import (
	"fmt"
	"testing"

	sq "github.com/Masterminds/squirrel"
	bunx "github.com/blink-io/hyperbun"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/stdscan"
	"github.com/stretchr/testify/require"
)

func TestSqlite_Builder_Select_1(t *testing.T) {
	db := getSqliteDB()

	b := bunx.B()
	type IdName struct {
		ID   string `db:"id"`
		GUID string `db:"guid"`
		Name string `db:"name"`
	}
	sql, args, err := b.Select("id", "name").
		From("applications").
		Where(sq.Gt{"id": 0}).
		Limit(5).
		ToSql()
	require.NoError(t, err)

	rts, err := stdscan.All(ctx, db, scan.StructMapper[*IdName](), sql, args...)
	require.NoError(t, err)

	fmt.Println(rts)
}
