package model

type IDModel struct {
	ID   int64  `bun:"id,pk,autoincrement" db:"id,pk"`
	GUID string `bun:"guid,unique,notnull" db:"guid"`
}

type IDSetter struct {
	ID   *int64  `db:"id"`
	GUID *string `db:"guid"`
}
