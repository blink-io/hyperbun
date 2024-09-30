package pg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBun_Pg_Raw_Select_Custom_1(t *testing.T) {
	db := getPgDB()
	var rs string
	err := db.NewRaw("select version()").
		Scan(context.Background(), &rs)
	require.NoError(t, err)
	fmt.Printf("pg version: %s\n", rs)
}

func TestBun_Pg_Raw_Select_Custom_2(t *testing.T) {
	db := getPgDB()
	var rs time.Time
	err := db.NewRaw("select created_at from tags where id = 5 limit 1").
		Scan(context.Background(), &rs)

	require.NoError(t, err)
	fmt.Printf("time: %s\n", rs)
}
