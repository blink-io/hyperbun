package sqlite

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	bunx "github.com/blink-io/hyperbun"
	"github.com/blink-io/hyperbun/extra/timing"

	sqlx "github.com/blink-io/hypersql"
	logginghook "github.com/blink-io/hypersql/driver/hooks/logging"
	"github.com/qustavo/sqlhooks/v2"
)

var ctx = context.Background()

var sqlitePath = filepath.Join(".", "sqlite_demo.db")

func dbOpts() []bunx.Option {
	opts := []bunx.Option{
		bunx.WithQueryHooks(timing.New()),
	}
	return opts
}

func sqliteCfg() *sqlx.Config {
	rpath, _ := filepath.Abs(sqlitePath)
	fmt.Println("Real path for sqlite: ", rpath)

	var cfg = &sqlx.Config{
		Dialect:     sqlx.DialectSQLite,
		Host:        sqlitePath,
		DriverHooks: newDriverHooks(),
		Logger: func(format string, args ...any) {
			msg := fmt.Sprintf(format, args...)
			slog.Default().With("db", "sqlite").Info(msg, "mode", "test")
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

func getSqliteFuncMap() map[string]string {
	funcsMap := map[string]string{
		"hex":              "hex(randomblob(32))",
		"random":           "random()",
		"version":          "sqlite_version()",
		"changes":          "changes()",
		"total_changes":    "total_changes()",
		"lower":            `lower("HELLO")`,
		"upper":            `upper("hello")`,
		"length":           `length("hello")`,
		"sqlite_source_id": `sqlite_source_id()`,
		//`concat("Hello", ",", "World")`,
	}
	return funcsMap
}
