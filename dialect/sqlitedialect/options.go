package sqlitedialect

import (
	"time"
)

type options struct {
	loc *time.Location
}

type Option func(o *options)

func applyOptions(ops ...Option) *options {
	opt := new(options)
	for _, o := range ops {
		o(opt)
	}
	return opt
}

// WithLocation uses time location
func WithLocation(loc *time.Location) Option {
	return func(o *options) {
		o.loc = loc
	}
}
