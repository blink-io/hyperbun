package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	bunx "github.com/blink-io/hyperbun"
	"github.com/blink-io/hyperbun/extra/timing"
	sqlx "github.com/blink-io/hypersql"
	logginghook "github.com/blink-io/hypersql/driver/hooks/logging"
	"github.com/qustavo/sqlhooks/v2"
)

var ctx = context.Background()

func dbOpts() []bunx.Option {
	opts := []bunx.Option{
		bunx.WithQueryHooks(timing.New()),
	}
	return opts
}

func pgCfg() *sqlx.Config {
	var cfg = &sqlx.Config{
		Dialect:     sqlx.DialectPostgres,
		Host:        "192.168.50.88",
		Port:        5432,
		User:        "test",
		Password:    "test",
		Name:        "test",
		DriverHooks: newDriverHooks(),
		Logger: func(format string, args ...any) {
			msg := fmt.Sprintf(format, args...)
			slog.Default().With("db", "postgres").Info(msg, "mode", "test")
		},
		Loc: time.Local,
	}
	return cfg
}

func newDriverHooks() []sqlhooks.Hooks {
	hs := []sqlhooks.Hooks{
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

func getPgDB() *bunx.DB {
	db, err := bunx.NewFromConf(pgCfg(), dbOpts()...)

	if err != nil {
		panic(err)
	}

	return db
}
