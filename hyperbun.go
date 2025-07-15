package hyperbun

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"

	"github.com/blink-io/hyperbun/errwrap"

	"github.com/uptrace/bun"
)

type DB interface {
	Context() context.Context
	NewSelect() *bun.SelectQuery
	NewInsert() *bun.InsertQuery
	NewUpdate() *bun.UpdateQuery
	NewDelete() *bun.DeleteQuery
	NewMerge() *bun.MergeQuery
	NewRaw(string, ...any) *bun.RawQuery
	NewValues(model any) *bun.ValuesQuery
	RunInTx(fn func(tx TxContext) error) error
	ForceRunInTx(fn func(tx TxContext) error) error
}

type IDType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~string
}

type Context struct {
	ctx context.Context
	Bun *bun.DB
}

var _ DB = (*Context)(nil)

func NewContext(ctx context.Context, db *bun.DB) *Context {
	return &Context{
		ctx: ctx,
		Bun: db,
	}
}

func (c Context) Context() context.Context {
	return c.ctx
}

func (c Context) NewSelect() *bun.SelectQuery {
	return c.Bun.NewSelect()
}

func (c Context) NewInsert() *bun.InsertQuery {
	return c.Bun.NewInsert()
}

func (c Context) NewUpdate() *bun.UpdateQuery {
	return c.Bun.NewUpdate()
}

func (c Context) NewDelete() *bun.DeleteQuery {
	return c.Bun.NewDelete()
}

func (c Context) NewMerge() *bun.MergeQuery {
	return c.Bun.NewMerge()
}

func (c Context) NewRaw(query string, args ...any) *bun.RawQuery {
	return c.Bun.NewRaw(query, args...)
}

func (c Context) NewValues(model any) *bun.ValuesQuery {
	return c.Bun.NewValues(model)
}

func (c Context) RunInTx(fn func(tx TxContext) error) error {
	return c.Bun.RunInTx(c.ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		return fn(NewTxContext(ctx, tx))
	})
}

func (c Context) ForceRunInTx(fn func(tx TxContext) error) error {
	return c.RunInTx(fn)
}

// ---------------------------------------------------------------

type TxContext struct {
	ctx context.Context
	Bun bun.Tx
}

var _ DB = (*TxContext)(nil)

func NewTxContext(ctx context.Context, tx bun.Tx) TxContext {
	return TxContext{
		ctx: ctx,
		Bun: tx,
	}
}

func (tx TxContext) Context() context.Context {
	return tx.ctx
}

func (tx TxContext) NewSelect() *bun.SelectQuery {
	return tx.Bun.NewSelect()
}

func (tx TxContext) NewInsert() *bun.InsertQuery {
	return tx.Bun.NewInsert()
}

func (tx TxContext) NewUpdate() *bun.UpdateQuery {
	return tx.Bun.NewUpdate()
}

func (tx TxContext) NewDelete() *bun.DeleteQuery {
	return tx.Bun.NewDelete()
}

func (tx TxContext) NewMerge() *bun.MergeQuery {
	return tx.Bun.NewMerge()
}

func (tx TxContext) NewRaw(query string, args ...any) *bun.RawQuery {
	return tx.Bun.NewRaw(query, args...)
}

func (tx TxContext) NewValues(model any) *bun.ValuesQuery {
	return tx.Bun.NewValues(model)
}

func (tx TxContext) RunInTx(fn func(tx TxContext) error) error {
	return fn(tx)
}

func (tx TxContext) ForceRunInTx(fn func(tx TxContext) error) error {
	return tx.Bun.RunInTx(tx.ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		return fn(NewTxContext(ctx, tx))
	})
}

func ByID[T any, ID IDType](db DB, id ID) (*T, error) {
	var row T
	if err := db.NewSelect().
		Model(&row).
		Where("id = ?", id).
		Limit(1).
		Scan(db.Context()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "ByID", "table", tableForType[T](), "id", id)
	}

	return &row, nil
}

func StructByID[T any, ID IDType](db DB, table string, id ID) (*T, error) {
	var row T
	columns := getColumns(reflect.TypeOf(row))
	if err := db.NewSelect().
		Column(columns...).
		Table(table).
		Where("id = ?", id).
		Limit(1).
		Scan(db.Context(), &row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "StructByID", "table", table, "id", id)
	}

	return &row, nil
}

func TypeByID[T any, ID IDType](db DB, table string, column string, id ID) (*T, error) {
	var value T
	if err := db.NewSelect().
		ColumnExpr(column).
		Table(table).
		Where("id = ?", id).
		Limit(1).
		Scan(db.Context(), &value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "TypeByID", "table", table, "column", column, "id", id)
	}

	return &value, nil
}

func BySQL[T any](db DB, query string, args ...any) (*T, error) {
	var row T
	if err := db.NewSelect().
		Model(&row).
		Where(query, args...).
		Limit(1).
		Scan(db.Context()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "BySQL", "table", tableForType[T]())
	}

	return &row, nil
}

func StructBySQL[T any](db DB, table string, query string, args ...any) (*T, error) {
	var row T
	columns := getColumns(reflect.TypeOf(row))
	if err := db.NewSelect().
		Column(columns...).
		Table(table).
		Where(query, args...).
		Limit(1).
		Scan(db.Context(), &row); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "StructBySQL", "table", table)
	}

	return &row, nil
}

func TypeBySQL[T any](db DB, table string, column string, query string, args ...any) (*T, error) {
	var value T
	if err := db.NewSelect().
		ColumnExpr(column).
		Table(table).
		Where(query, args...).
		Limit(1).
		Scan(db.Context(), &value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "TypeBySQL", "table", table, "column", column)
	}

	return &value, nil
}

func Many[T any](db DB, query string, args ...any) ([]T, error) {
	var rows []T
	if err := db.NewSelect().
		Model(&rows).
		Where(query, args...).
		Scan(db.Context()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, annotate(err, "Many", "table", tableForType[T]())
	}

	return rows, nil
}

func Exists[ID IDType](db DB, table string, id ID) (bool, error) {
	c, err := CountQuery(db, table, "id = ?", id)
	if err != nil {
		return false, annotate(err, "Exists", "table", table, "id", id)
	}

	return c == 1, nil
}

func ExistsBySQL(db DB, table string, query string, args ...any) (bool, error) {
	var exists bool
	if err := db.NewRaw("SELECT EXISTS(SELECT 1 from "+table+" WHERE "+query+")", args...).
		Scan(db.Context(), &exists); err != nil {
		return false, annotate(err, "ExistsBySQL", "table", table)
	}

	return exists, nil
}

func CountQuery(db DB, table string, query string, args ...any) (int, error) {
	count, err := db.NewSelect().
		Table(table).
		Where(query, args...).
		Count(db.Context())
	if err != nil {
		return 0, annotate(err, "CountQuery", "table", table)
	}
	return count, nil
}

func Insert[T any](db DB, row *T) error {
	_, err := db.NewInsert().
		Model(row).
		Exec(db.Context())
	if err != nil {
		return annotate(err, "Insert", "table", tableForType[T]())
	}
	return nil
}

func InsertMany[T any](db DB, rows []T) error {
	if len(rows) == 0 {
		return nil
	}

	_, err := db.NewInsert().
		Model(&rows).
		Exec(db.Context())
	if err != nil {
		return annotate(err, "InsertMany", "table", tableForType[T]())
	}
	return nil
}

func Update[T any](db DB, row *T, pk ...string) error {
	if len(pk) == 0 {
		pk = append(pk, "id")
	}
	if _, err := db.NewUpdate().
		Model(row).
		WherePK(pk...).
		Exec(db.Context()); err != nil {
		return annotate(err, "Update", "table", tableForType[T](), "pk", strings.Join(pk, ","))
	}
	return nil
}

func UpdateSQLByID[ID IDType](db DB, table string, id ID, query string, args ...any) error {
	_, err := db.NewUpdate().
		Table(table).
		Set(query, args...).
		Where("id = ?", id).
		Exec(db.Context())
	if err != nil {
		return annotate(err, "UpdateSQLByID", "table", table, "id", id)
	}
	return err
}

// Upsert defines upsert and check multiple constraints, see
// https://stackoverflow.com/questions/35888012/use-multiple-conflict-target-in-on-conflict-clause
func Upsert[T any](db DB, rows T, conflictColumns string) error {
	if _, err := db.NewInsert().
		Model(rows).
		On(fmt.Sprintf("conflict (%s) do update", conflictColumns)).
		Exec(db.Context()); err != nil {
		return annotate(err, "Upsert", "table", tableForType[T](), "conflict", conflictColumns)
	}
	return nil
}

func UpsertIgnore[T any](db DB, rows T) error {
	_, err := db.NewInsert().
		Model(rows).
		On("conflict do nothing").
		Exec(db.Context())
	if err != nil {
		return annotate(err, "UpsertIgnore", "table", tableForType[T]())
	}

	return err
}

func DeleteByID[ID IDType](db DB, table string, id ID) error {
	if _, err := db.NewDelete().
		Table(table).
		Where("id = ?", id).
		Exec(db.Context()); err != nil {
		return annotate(err, "DeleteByID", "table", table, "id", id)
	}

	return nil
}

func DeleteBySQL(db DB, table string, query string, args ...any) error {
	if _, err := db.NewDelete().
		Table(table).
		Where(query, args...).
		Exec(db.Context()); err != nil {
		return annotate(err, "DeleteBySQL", "table", table)
	}

	return nil
}

func RunInTx(db DB, fn func(tx TxContext) error) error {
	if err := db.RunInTx(fn); err != nil {
		return fmt.Errorf("RunInTx: %w", err)
	}
	return nil
}

func ForceRunInTx(db DB, fn func(tx TxContext) error) error {
	if err := db.ForceRunInTx(fn); err != nil {
		return fmt.Errorf("ForceRunInTx: %w", err)
	}
	return nil
}

func RunInLockedTx(db DB, id string, fn func(tx TxContext) error) error {
	return RunInTx(db, func(tx TxContext) error {
		if err := advisoryLock(db, id); err != nil {
			return errwrap.Wrap(err)
		}

		return fn(tx)
	})
}

func advisoryLock(db DB, name string) error {
	h := fnv.New64()
	if _, err := h.Write([]byte(name)); err != nil {
		return errwrap.Wrap(err)
	}
	s := h.Sum64()
	if _, err := db.NewRaw("SELECT pg_advisory_xact_lock(?)", int64(s)).
		Exec(db.Context()); err != nil {
		return errwrap.Wrap(err)
	}

	return nil
}

func annotate(err error, op string, kvs ...any) error {
	pairs := make([][2]string, len(kvs)/2)
	numPairs := len(kvs) / 2
	odd := len(kvs)%2 != 0
	for i := 0; i < numPairs; i++ {
		pairs[i] = [2]string{
			fmt.Sprint(kvs[i*2]),
			fmt.Sprint(kvs[i*2+1]),
		}
	}
	if odd {
		pairs = append(pairs, [2]string{
			fmt.Sprint(kvs[len(kvs)-1]),
			"<missing value>",
		})
	}
	joined := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		joined = append(joined, fmt.Sprint(pair[0], "='", pair[1], "'"))
	}
	joinedStr := strings.Join(joined, " ")
	if joinedStr != "" {
		joinedStr = " " + joinedStr
	}

	return fmt.Errorf("performing %s%s: %w", op, joinedStr, err)
}

func tableForType[T any]() string {
	var t T
	typ := reflect.TypeOf(t)
	kind := typ.Kind()

	// This covers the case like *U, []U, *[]U, []*U etc
	for kind == reflect.Pointer || kind == reflect.Slice {
		typ = typ.Elem()
		kind = typ.Kind()
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		val, ok := f.Tag.Lookup("bun")
		if !ok {
			continue
		}
		for _, ann := range strings.Split(val, ",") {
			spl := strings.Split(ann, ":")
			if len(spl) != 2 {
				continue
			}
			if spl[0] == "table" {
				return spl[1]
			}
		}
	}
	return ""
}
