package model

import "github.com/blink-io/opt/omit"

type IDModel struct {
	ID   int64  `bun:"id,pk,nullzero,autoincrement" db:"id,pk"`
	GUID string `bun:"guid,unique,notnull,type:varchar(60)" db:"guid"`
}

type IDSetter struct {
	ID   omit.Val[int64]  `db:"id"`
	GUID omit.Val[string] `db:"guid"`
}
