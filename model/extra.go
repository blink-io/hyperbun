package model

import (
	"time"

	"github.com/aarondl/opt/omitnull"
)

type ExtraModel struct {
	CreatedAt time.Time `bun:"created_at,notnull,skipupdate" db:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull" db:"updated_at"`
	// Optional fields for tables
	CreatedBy omitnull.Val[string]    `bun:"created_by,type:varchar(60),nullzero,skipupdate" db:"created_by"`
	UpdatedBy omitnull.Val[string]    `bun:"updated_by,type:varchar(60),nullzero" db:"updated_by"`
	DeletedAt omitnull.Val[time.Time] `bun:"deleted_at,nullzero,skipupdate" db:"deleted_at"`
	DeletedBy omitnull.Val[string]    `bun:"deleted_by,type:varchar(60),nullzero,skipupdate" db:"deleted_by"`
	IsDeleted omitnull.Val[bool]      `bun:"is_deleted,nullzero,skipupdate" db:"is_deleted"`
}
