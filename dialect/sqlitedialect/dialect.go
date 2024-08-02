package sqlitedialect

import (
	"time"

	"github.com/uptrace/bun/dialect/sqlitedialect"
)

type Dialect struct {
	*sqlitedialect.Dialect
	opts *options
}

func New(ops ...Option) *Dialect {
	opts := applyOptions(ops...)
	d := &Dialect{Dialect: sqlitedialect.New(), opts: opts}
	d.Features()
	return d
}

// AppendTime in the schema.BaseDialect uses the UTC timezone.
// Let the developers make the decision.
func (d *Dialect) AppendTime(b []byte, tm time.Time) []byte {
	b = append(b, '\'')
	if loc := d.opts.loc; loc != nil {
		tm = tm.In(loc)
	}
	b = tm.AppendFormat(b, time.RFC3339Nano)
	b = append(b, '\'')
	return b
}
