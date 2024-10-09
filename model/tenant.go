package model

import "github.com/blink-io/opt/omit"

type TenantModel struct {
	TenantID int `bun:"tenant_id,notnull" db:"tenant_id,pk" json:"tenant_id,omitempty" toml:"tenant_id,omitempty" yaml:"tenant_id,omitempty" msgpack:"tenant_id,omitempty"`
}

type TenantSetter struct {
	TenantID omit.Val[int64] `db:"tenant_id"`
}
