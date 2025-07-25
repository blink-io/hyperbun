package sqlite

import (
	"time"

	"github.com/blink-io/opt/omitnull"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sanity-io/litter"
	"github.com/uptrace/bun"
)

// Model8 for dbx
type Model8 struct {
	ID      int    `db:"pk"`
	Name    string `db:"name"`
	Title   string `db:"title"`
	Fax     string `db:"fax"`
	Web     string `db:"web"`
	Age     int    `db:"age"`
	Right   bool   `db:"right"`
	Counter int64  `db:"counter"`
}

func NewModel8() *Model8 {
	m := new(Model8)
	m.Name = "Orm Benchmark"
	m.Title = "Just a Benchmark for fun"
	m.Fax = "99909990"
	m.Web = "http://blog.milkpod29.me"
	m.Age = 100
	m.Right = true
	m.Counter = 1000

	return m
}

func (Model8) Table() string {
	return "models"
}

func (Model8) TableName() string {
	return "models"
}

type UserWithDevices struct {
	bun.BaseModel `bun:"users,alias:u_1" db:"-" json:"-" toml:"-" yaml:"-" msgpack:"-"`
	ID            int64  `bun:"id,pk" db:"id"`
	Username      string `bun:"username,type:varchar(60),notnull" db:"username"`

	Devices []*UserDevice `bun:"ref:has-many, join:id=user_id"`
}

// User represents iOS/Android/Windows/OSX/Linux application
type User struct {
	bun.BaseModel `bun:"users,alias:u_1" db:"-" json:"-" toml:"-" yaml:"-" msgpack:"-"`
	ID            int64                `bun:"id,pk,notnull" db:"id"`
	Username      string               `bun:"username,type:varchar(60),notnull" db:"username"`
	Location      string               `bun:"location,type:varchar(60),notnull" db:"location"`
	Profile       string               `bun:"profile,type:varchar(200),notnull" db:"profile"`
	Level         int8                 `bun:"level,notnull" db:"code" json:"level,omitempty"`
	Description   omitnull.Val[string] `bun:"description,type:text" db:"description"`
}

func (User) TableName() string {
	return "users"
}

func (User) Table() string {
	return "users"
}

// UserDevice represents user's devices
type UserDevice struct {
	bun.BaseModel `bun:"user_devices,alias:ud_1" db:"-" json:"-" toml:"-" yaml:"-" msgpack:"-"`
	ID            int64                `bun:"id,pk,notnull" db:"id"`
	Name          string               `bun:"name,type:varchar(60),notnull" db:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty" msgpack:"name,omitempty"`
	Description   omitnull.Val[string] `bun:"description,type:text" db:"description" json:"description,omitempty" toml:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
	CreatedAt     time.Time            `bun:"created_at,notnull,skipupdate" db:"created_at" json:"created_at,omitempty" toml:"created_at,omitempty" yaml:"created_at,omitempty" msgpack:"created_at,omitempty"`
	UpdatedAt     time.Time            `bun:"updated_at,notnull" db:"updated_at" json:"updated_at,omitempty" toml:"updated_at,omitempty" yaml:"updated_at,omitempty" msgpack:"updated_at,omitempty"`
}

func (UserDevice) TableName() string {
	return "user_devices"
}

func (UserDevice) Table() string {
	return "user_devices"
}

// Application represents iOS/Android/Windows/OSX/Linux application
type Application struct {
	bun.BaseModel `bun:"applications,alias:applications" db:"-" json:"-" toml:"-" yaml:"-" msgpack:"-"`
	ID            int64                `bun:"id,pk,notnull" db:"id"`
	Level         int32                `bun:"level,type:integer,notnull" db:"level" json:"level,omitempty" toml:"level,omitempty" yaml:"level,omitempty" msgpack:"level,omitempty"`
	Status        string               `bun:"status,type:varchar(60),notnull" db:"status" json:"status,omitempty" toml:"status,omitempty" yaml:"status,omitempty" msgpack:"status,omitempty"`
	Type          string               `bun:"type,type:varchar(60),notnull" db:"type" json:"type,omitempty" toml:"type,omitempty" yaml:"type,omitempty" msgpack:"type,omitempty"`
	Name          string               `bun:"name,type:varchar(200),notnull" db:"name" json:"name,omitempty" toml:"name,omitempty" yaml:"name,omitempty" msgpack:"name,omitempty"`
	Code          string               `bun:"code,type:varchar(60),unique,notnull" db:"code" json:"code,omitempty" toml:"code,omitempty" yaml:"code,omitempty" msgpack:"code,omitempty"`
	Description   omitnull.Val[string] `bun:"description,type:text" db:"description" json:"description,omitempty" toml:"description,omitempty" yaml:"description,omitempty" msgpack:"description,omitempty"`
}

func (Application) TableName() string {
	return "applications"
}

func (Application) Table() string {
	return "applications"
}

func (Application) Columns() []string {
	return appColumns
}

var appColumns = []string{
	"id",
	"guid",
	"status",
	"name",
	"code",
	"type",
	"created_at",
	"updated_at",
	"deleted_at",
}

func printApp(r *Application) string {
	return litter.Sdump(r)
}

func ToAnySlice[T any](a []T) []any {
	t := make([]any, 0, len(a))
	for _, v := range a {
		t = append(t, v)
	}
	return t
}

type IDAndName struct {
	// BaseModel is needed for applying table name
	bun.BaseModel `bun:"table:applications,alias:a1"`
	ID            int64  `bun:"id,type:bigint,pk"`
	Name          string `bun:"name,type:text"`
}

type IDAndProfile struct {
	//bun.BaseModel `bun:"table:applications,alias:a1"`
	ID      int64  `bun:"id,type:bigint,pk"`
	GUID    string `bun:"guid,type:text"`
	Profile string `bun:"profile,type:text"`
}

func randomApplication() *Application {
	r := new(Application)
	r.Name = gofakeit.Name()
	r.Code = gofakeit.UUID()
	r.Level = int32(gofakeit.IntRange(1, 99))
	r.Status = "ok"
	r.Type = gofakeit.FileMimeType()
	r.Description = omitnull.From(gofakeit.ProductDescription())
	return r
}
