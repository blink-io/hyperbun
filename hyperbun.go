package hyperbun

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"reflect"

	"github.com/blink-io/hypersql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

var (
	ErrNilConfig = errors.New("[hyperbun] config is nil")
)

type Config = hypersql.Config

type (
	ext interface {
		RawDB() *RawDB

		RegisterModel(m ...any)

		Table(typ reflect.Type) *schema.Table
	}

	IDB interface {
		RawIDB

		io.Closer

		ext
	}

	rdb = bun.DB

	DB struct {
		*rdb
		sqlDB *sql.DB
	}
)

var _ IDB = (*DB)(nil)

func NewFromSqlDB(sqlDB *sql.DB, dialect Dialect, ops ...Option) (*DB, error) {
	rdb := bun.NewDB(sqlDB, dialect, bun.WithDiscardUnknownColumns())

	opts := applyOptions(ops...)
	for _, h := range opts.queryHooks {
		rdb.AddQueryHook(h)
	}

	db := &DB{
		rdb:   rdb,
		sqlDB: sqlDB,
	}

	return db, nil
}

func NewFromConf(c *Config, ops ...Option) (*DB, error) {
	if c == nil {
		return nil, ErrNilConfig
	}

	ctx := context.Background()
	dOpts := make([]DialectOption, 0)
	if c.Loc != nil {
		dOpts = append(dOpts, DialectWithLoc(c.Loc))
	}
	dialect, err := GetDialect(ctx, c, dOpts...)
	if err != nil {
		return nil, err
	}

	sqlDB, err := hypersql.NewSqlDB(c)
	if err != nil {
		return nil, err
	}

	return NewFromSqlDB(sqlDB, dialect, ops...)
}

func (db *DB) RegisterModel(m ...any) {
	db.rdb.RegisterModel(m...)
}

func (db *DB) SqlDB() *sql.DB {
	return db.rdb.DB
}

func (db *DB) Close() error {
	if db.rdb != nil {
		return db.rdb.Close()
	}
	return nil
}

func (db *DB) RawDB() *RawDB {
	return db.rdb
}

func (db *DB) HealthCheck(ctx context.Context) error {
	return hypersql.DoPingContext(ctx, db.sqlDB)
}
