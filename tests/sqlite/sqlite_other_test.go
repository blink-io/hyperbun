package sqlite

import (
	"testing"

	xsql "github.com/blink-io/hypersql"
	"github.com/stretchr/testify/require"
)

func getXSqlOpts(d, n string) *xsql.Config {
	opts := &xsql.Config{
		Dialect: d,
		Name:    n,
	}
	return opts
}

type wrapOptions struct {
	cfg *xsql.Config
}

func TestSqlOptions(t *testing.T) {
	opts1 := getXSqlOpts("ddd", "nnnn")
	opts2 := getXSqlOpts("ffff", "aaaa")
	opts3 := getXSqlOpts("qqqq", "tttt")

	as := []*xsql.Config{
		opts1,
		opts2,
	}

	as = append(as, opts3)

	wopt := wrapOptions{
		cfg: opts1,
	}

	opts1.Dialect = "ninin"
	opts1.Name = "kkdf"

	opts3.Dialect = "gggg"
	opts3.Name = "vvvv"

	require.NotNil(t, opts1)
	require.NotNil(t, wopt)
	require.NotNil(t, as)
}
