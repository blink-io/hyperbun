package schema

import (
	"fmt"
	"github.com/sanity-io/litter"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	bschema "github.com/uptrace/bun/schema"
	"reflect"
	"testing"
)

func TestSchema_1(t *testing.T) {
	type MyTable struct {
		Name string
		Len  int
	}
	type MyTableColumns struct {
		Name  Column
		Level Column
	}
	var myTableColumns = MyTableColumns{
		Name:  "name",
		Level: "level",
	}

	var schema = Schema[MyTable, MyTableColumns]{
		PrimaryKeys: PrimaryKeys{"id1", "id2"},
		Table:       "MyTable",
		Label:       "MyTable",
		Model:       (*MyTable)(nil),
		Alias:       "MyTable",
		Columns:     myTableColumns,
	}

	fmt.Println(litter.Sdump(schema))
	fmt.Printf("Type of Model: %#v\n", reflect.TypeOf(schema.Model).String())

	var nn = schema.Columns.Level
	var tt = schema.Table
	fmt.Println(nn)
	fmt.Println(tt)

	var fmter1 = bschema.NewFormatter(pgdialect.New())
	var fmter2 = bschema.NewFormatter(mysqldialect.New())

	var bb1 []byte
	dd1, err1 := bun.Name(schema.Columns.Level).AppendQuery(fmter1, bb1)
	require.NoError(t, err1)
	fmt.Println(string(dd1))

	var bb2 []byte
	dd2, err2 := bun.Name(schema.Columns.Name).AppendQuery(fmter2, bb2)
	require.NoError(t, err2)
	fmt.Println(string(dd2))
}
