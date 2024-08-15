package hyperbun

import (
	sq "github.com/Masterminds/squirrel"
)

type Builder = sq.StatementBuilderType

var sb = sq.StatementBuilder

func (db *DB) B() Builder {
	return sb
}

func B() Builder {
	return sb
}
