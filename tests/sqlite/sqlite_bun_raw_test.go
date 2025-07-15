package sqlite

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSqlite_Bun_Raw_Select_Custom_2(t *testing.T) {
	db := getSqliteBunDB()
	var rs string
	err := db.NewRaw("select sqlite_version()").
		Scan(ctx, &rs)
	require.NoError(t, err)
	fmt.Printf("sqlite version: %s\n", rs)
}
