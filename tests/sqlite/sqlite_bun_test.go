package sqlite

import (
	"fmt"
	"testing"

	"github.com/blink-io/opt/omitnull"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

var m1 = (*Application)(nil)
var m2 = (*User)(nil)
var ms = []any{m1, m2}

func TestSqlite_Bun_CreateTable_1(t *testing.T) {
	db := getSqliteBunDB()
	for _, m := range ms {
		_, err := db.NewCreateTable().IfNotExists().Model(m).Exec(ctx)
		require.NoError(t, err)
	}
}

func TestSqlite_Bun_DropTable_1(t *testing.T) {
	db := getSqliteBunDB()
	for _, m := range ms {
		_, err := db.NewDropTable().IfExists().Model(m).Exec(ctx)
		require.NoError(t, err)
	}
}

func TestSqlite_Bun_RebuildTable_1(t *testing.T) {
	TestSqlite_Bun_DropTable_1(t)
	TestSqlite_Bun_CreateTable_1(t)
}

func TestSqlite_Bun_Custom_Update_1(t *testing.T) {
	db := getSqliteBunDB()

	ds := db.NewUpdate().
		Table("applications").
		SetColumn("status", "?", "no-ok").
		Where("status = ?", "ok")
	_, err1 := ds.Exec(ctx)
	require.NoError(t, err1)
}

func TestSqlite_Bun_Custom_Update_Bulk_1(t *testing.T) {
	db := getSqliteBunDB()

	u1 := new(User)
	u1.ID = 1
	u1.Location = gofakeit.City()
	u1.Profile = gofakeit.AppName()

	u2 := new(User)
	u2.ID = 2
	u2.Location = gofakeit.City()
	u2.Profile = gofakeit.AppName()

	values := db.NewValues(&[]*User{u1, u2})

	_, err := db.NewUpdate().
		With("_data", values).
		Model((*User)(nil)).
		TableExpr("_data").
		OmitZero().
		Set("location = _data.location").
		Set("profile = _data.profile").
		Where("users.id = _data.id").
		Exec(ctx)

	require.NoError(t, err)
}

func TestSqlite_Bun_Custom_Update_Bulk_2(t *testing.T) {
	db := getSqliteBunDB()

	u1 := map[string]any{
		"id":       1,
		"location": gofakeit.City() + "DoUpdate-Bulk-By-Map",
		"profile":  "Profile:" + "DoUpdate-Bulk-By-Map",
	}

	u2 := map[string]any{
		"id":       2,
		"location": gofakeit.City() + "DoUpdate-Bulk-By-Map",
		"profile":  "Profile:" + "DoUpdate-Bulk-By-Map",
	}
	values := db.NewValues(&[]map[string]any{u1, u2})

	_, err := db.NewUpdate().
		With("_data", values).
		Model((*User)(nil)).
		TableExpr("_data").
		OmitZero().
		Set("location = _data.location").
		Set("profile = _data.profile").
		Where("users.id = _data.id").
		Exec(ctx)

	require.NoError(t, err)
}

func TestSqlite_Bun_Map_Update_1(t *testing.T) {
	db := getSqliteBunDB()

	value := map[string]interface{}{
		"status":      "no-ok",
		"description": "updated-by-map",
	}
	_, err := db.NewUpdate().
		Model(&value).
		TableExpr("applications").
		Where("id = ?", 20).
		Exec(ctx)
	require.NoError(t, err)
}

func TestSqlite_Bun_RawSQL_Update_1(t *testing.T) {
	db := getSqliteBunDB()

	_, err := db.NewUpdate().NewRaw(
		"update applications set status = ?, description = ? where id = ?",
		"raw-sql-ok", "update-by-raw-sql", 22,
	).Exec(ctx)
	require.NoError(t, err)
}

func TestSqlite_Bun_Map_Insert_1(t *testing.T) {
	db := getSqliteBunDB()

	val := randomApplication()

	sqlstr := db.NewInsert().
		Model(&val).
		Ignore().
		TableExpr("users").
		String()
	fmt.Println("sql: ", sqlstr)

	_, err := db.Exec(sqlstr)
	require.NoError(t, err)
}

func TestSqlite_Bun_Map_Insert_2(t *testing.T) {
	db := getSqliteBunDB()

	vals := []*Application{
		randomApplication(),
		randomApplication(),
	}

	q := db.NewInsert().
		Model(&vals).
		Ignore().
		Table("users")

	_, err := q.Exec(ctx)
	require.NoError(t, err)
}

func TestSqlite_Bun_Custom_Select_1(t *testing.T) {
	db := getSqliteBunDB()

	var ids []int64
	var descs []omitnull.Val[string]
	err := db.NewSelect().
		Table("applications").
		Column("id", "description").
		Limit(5).
		Scan(ctx, &ids, &descs)
	require.NoError(t, err)
}

func TestSqlite_Bun_Select_Funcs(t *testing.T) {
	db := getSqliteBunDB()

	type Result struct {
		Payload string `db:"payload"`
	}

	sqlF := "select %s as payload"
	funcs := getSqliteFuncMap()
	for k, v := range funcs {
		ss := fmt.Sprintf(sqlF, v)
		row := db.QueryRowContext(ctx, ss)
		var rstr string
		err := row.Scan(&rstr)
		require.NoError(t, err)
		fmt.Println(k, "=>", rstr)
	}
}
