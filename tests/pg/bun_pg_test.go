package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/blink-io/hyperbun/extra/verbose"
	sqlx "github.com/blink-io/hypersql"
	logginghook "github.com/blink-io/hypersql/driver/hooks/logging"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var ctx = context.Background()

func pgCfg() *sqlx.Config {
	var cfg = &sqlx.Config{
		Dialect:     sqlx.DialectPostgres,
		Host:        "192.168.50.88",
		Port:        5432,
		User:        "test",
		Password:    "test",
		Name:        "test",
		DriverHooks: pgDriverHooks(),
		Logger: func(format string, args ...any) {
			msg := fmt.Sprintf(format, args...)
			slog.Default().With("db", "postgres").Info(msg, "mode", "test")
		},
		Loc: time.Local,
	}
	return cfg
}

func pgDriverHooks() sqlx.DriverHooks {
	hs := sqlx.DriverHooks{
		logginghook.Func(func(format string, args ...any) {
			slog.Default().Info(fmt.Sprintf(format, args...))
		}),
	}
	return hs
}
func getPgSqlDB() *sql.DB {
	db, err := sqlx.NewSqlDB(pgCfg())
	//db.AddQueryHook(logging.Func(log.Printf))
	if err != nil {
		panic(err)
	}

	return db
}

func getPgDB() *bun.DB {
	db := bun.NewDB(getPgSqlDB(), pgdialect.New(), bun.WithDiscardUnknownColumns())
	db.AddQueryHook(verbose.Default())
	return db
}
