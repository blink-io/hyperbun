package sqlite

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	sqlx "github.com/blink-io/hypersql"
	logginghook "github.com/blink-io/hypersql/driver/hooks/logging"
)

var ctx = context.Background()

var sqlitePath = filepath.Join(".", "sqlite_demo.db")

func sqliteCfg() *sqlx.Config {
	rpath, _ := filepath.Abs(sqlitePath)
	fmt.Println("Real path for sqlite: ", rpath)

	var cfg = &sqlx.Config{
		Dialect:     sqlx.DialectSQLite,
		Host:        sqlitePath,
		DriverHooks: sqliteDriverHooks(),
		Logger: func(format string, args ...any) {
			msg := fmt.Sprintf(format, args...)
			slog.Default().With("db", "sqlite").Info(msg, "mode", "test")
		},
		Loc: time.Local,
	}
	return cfg
}

func sqliteDriverHooks() sqlx.DriverHooks {
	hs := sqlx.DriverHooks{
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
