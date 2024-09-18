package model

import (
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
)

type ExtraModel struct {
	CreatedAt time.Time `bun:"created_at,notnull,skipupdate" db:"created_at" json:"created_at,omitempty" toml:"created_at,omitempty" yaml:"created_at,omitempty" msgpack:"created_at,omitempty"`
	UpdatedAt time.Time `bun:"updated_at,notnull" db:"updated_at" json:"updated_at,omitempty" toml:"updated_at,omitempty" yaml:"updated_at,omitempty" msgpack:"updated_at,omitempty"`
	// Optional fields for tables
	CreatedBy omitnull.Val[string]    `bun:"created_by,type:varchar(60),nullzero,skipupdate" db:"created_by" json:"created_by,omitempty" toml:"created_by,omitempty" yaml:"created_by,omitempty" msgpack:"created_by,omitempty"`
	UpdatedBy omitnull.Val[string]    `bun:"updated_by,type:varchar(60),nullzero" db:"updated_by" json:"updated_by,omitempty" toml:"updated_by,omitempty" yaml:"updated_by,omitempty" msgpack:"updated_by,omitempty"`
	DeletedAt omitnull.Val[time.Time] `bun:"deleted_at,nullzero,skipupdate" db:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at,omitempty" yaml:"deleted_at,omitempty" msgpack:"deleted_at,omitempty"`
	DeletedBy omitnull.Val[string]    `bun:"deleted_by,type:varchar(60),nullzero,skipupdate" db:"deleted_by" json:"deleted_by,omitempty" toml:"deleted_by,omitempty" yaml:"deleted_by,omitempty" msgpack:"deleted_by,omitempty"`
	IsDeleted omitnull.Val[bool]      `bun:"is_deleted,nullzero,skipupdate" db:"is_deleted" json:"is_deleted,omitempty" toml:"is_deleted,omitempty" yaml:"is_deleted,omitempty" msgpack:"is_deleted,omitempty"`
}

type ExtraSetter struct {
	CreatedAt omit.Val[time.Time] `db:"created_at"`
	UpdatedAt omit.Val[time.Time] `db:"updated_at"`
	// Optional fields for tables
	CreatedBy omitnull.Val[string]    `db:"created_by"`
	UpdatedBy omitnull.Val[string]    `db:"updated_by"`
	DeletedAt omitnull.Val[time.Time] `db:"deleted_at"`
	DeletedBy omitnull.Val[string]    `db:"deleted_by"`
	IsDeleted omitnull.Val[bool]      `db:"is_deleted"`
}
