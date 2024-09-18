package model

import "github.com/aarondl/opt/omit"

type IDModel struct {
	ID   int64  `bun:"id,pk,nullzero,autoincrement" db:"id,pk" json:"id,omitempty" toml:"id,omitempty" yaml:"id,omitempty" msgpack:"id,omitempty"`
	GUID string `bun:"guid,unique,notnull,type:varchar(60)" db:"guid" json:"guid,omitempty" toml:"guid,omitempty" yaml:"guid,omitempty" msgpack:"guid,omitempty"`
}

type IDSetter struct {
	ID   omit.Val[int64]  `db:"id"`
	GUID omit.Val[string] `db:"guid"`
}
