package sqlite

import (
	"database/sql"

	bunx "github.com/blink-io/hyperbun"
	sqlx "github.com/blink-io/hypersql"
)

func init() {
}

func getSqliteSqlDB() *sql.DB {
	db, err := sqlx.NewSqlDB(sqliteCfg())
	//db.AddQueryHook(logging.Func(log.Printf))
	if err != nil {
		panic(err)
	}

	return db
}

func getSqliteDB() *bunx.DB {
	db, err := bunx.NewFromConf(sqliteCfg(), dbOpts()...)

	if err != nil {
		panic(err)
	}

	return db
}
