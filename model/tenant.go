package model

type TenantModel struct {
	TenantID int64 `bun:"tenant_id,notnull" db:"tenant_id,pk" json:"tenant_id,omitempty" toml:"tenant_id,omitempty" yaml:"tenant_id,omitempty" msgpack:"tenant_id,omitempty"`
}
