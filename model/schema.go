package model

import (
	bunx "github.com/blink-io/hyperbun"
)

type Schema[M any, C any] struct {
	PK      string
	Label   string
	Alias   string
	Table   Table
	Model   *M
	Columns C
}

type Table string

func (t Table) String() string {
	return string(t)
}

type Column string

func (c Column) Name() bunx.Name {
	return bunx.Name(c)
}

func (c Column) Ident() bunx.Ident {
	return bunx.Ident(c)
}

func (c Column) Safe() bunx.Safe {
	return bunx.Safe(c)
}

func (c Column) String() string {
	return string(c)
}
