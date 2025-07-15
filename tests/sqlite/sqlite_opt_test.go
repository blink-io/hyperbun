package sqlite

import (
	"database/sql"

	"github.com/blink-io/hyperbun/extra/verbose"
	sqlx "github.com/blink-io/hypersql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

func init() {
}

func getSqliteSqlDB() *sql.DB {
	db, err := sqlx.NewSqlDB(sqliteCfg())
	if err != nil {
		panic(err)
	}

	return db
}

func getSqliteBunDB() *bun.DB {
	db := bun.NewDB(
		getSqliteSqlDB(),
		sqlitedialect.New(),
		bun.WithDiscardUnknownColumns(),
	)
	db.AddQueryHook(verbose.Default())
	return db
}
