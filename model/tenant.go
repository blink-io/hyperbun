package model

import "github.com/blink-io/opt/omit"

type TenantModel struct {
	TenantID int `bun:"tenant_id,notnull" db:"tenant_id"`
}

type TenantSetter struct {
	TenantID omit.Val[int64] `db:"tenant_id"`
}
