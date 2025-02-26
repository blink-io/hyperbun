package model

type IDModel struct {
	ID   int64  `bun:"id,pk,nullzero,autoincrement" db:"id,pk"`
	GUID string `bun:"guid,unique,notnull,type:varchar(60)" db:"guid"`
}
