package tests

import (
	"fmt"
	"testing"

	"github.com/blink-io/hyperbun/model"
)

func TestModel_1(t *testing.T) {
	cols := model.ColumnNames

	fmt.Println(cols.All())
}
