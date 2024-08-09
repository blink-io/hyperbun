package hyperbun

import "github.com/uptrace/bun"

type Logf func(format string, v ...any)

func (l Logf) Printf(format string, v ...interface{}) {
	l(format, v...)
}

func SetLogger(f Logf) {
	bun.SetLogger(f)
}
