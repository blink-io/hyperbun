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
		Name  string
		Level string
	}
	var myTableColumns = MyTableColumns{
		Name:  "name",
		Level: "level",
	}

	var table = Table[MyTable, MyTableColumns]{
		PrimaryKeys: []string{"id1", "id2"},
		Schema:      "public",
		Name:        "MyTable",
		Model:       (*MyTable)(nil),
		Alias:       "MyTable",
		Columns:     myTableColumns,
	}

	fmt.Println(litter.Sdump(table))
	fmt.Printf("Type of Model: %#v\n", reflect.TypeOf(table.Model).String())

	var nn = table.Columns.Level
	var tt = table.Name
	fmt.Println(nn)
	fmt.Println(tt)

	var fmter1 = bschema.NewFormatter(pgdialect.New())
	var fmter2 = bschema.NewFormatter(mysqldialect.New())

	var bb1 []byte
	dd1, err1 := bun.Name(table.Columns.Level).AppendQuery(fmter1, bb1)
	require.NoError(t, err1)
	fmt.Println(string(dd1))

	var bb2 []byte
	dd2, err2 := bun.Name(table.Columns.Name).AppendQuery(fmter2, bb2)
	require.NoError(t, err2)
	fmt.Println(string(dd2))
}
