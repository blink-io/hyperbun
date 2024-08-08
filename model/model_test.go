package model

import (
	"fmt"
	"testing"
)

func TestModel_1(t *testing.T) {
	cols := ColumnNames

	fmt.Println(cols.ID)
	fmt.Println(cols.GUID)
	fmt.Println(cols.CreatedAt)
	fmt.Println(cols.UpdatedAt)
}
