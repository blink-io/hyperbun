package model

import (
	"time"

	"github.com/aarondl/opt/omitnull"
)

type ExtraModel struct {
	CreatedAt time.Time `bun:"created_at,notnull,skipupdate" db:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull" db:"updated_at"`
	CreatedBy string    `bun:"created_by,notnull,skipupdate" db:"created_by"`
	UpdatedBy string    `bun:"updated_by,notnull" db:"updated_by"`
	// Optional fields for tables
	DeletedAt omitnull.Val[time.Time] `bun:"deleted_at,skipupdate" db:"deleted_at"`
	DeletedBy omitnull.Val[string]    `bun:"deleted_by,skipupdate" db:"deleted_by"`
	IsDeleted bool                    `bun:"is_deleted,skipupdate" db:"is_deleted"`
}
