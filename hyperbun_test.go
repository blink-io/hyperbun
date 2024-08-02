package hyperbun

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeQuery_1(t *testing.T) {
	q := "abc"
	ss := doSafeQuery(q, "a", "b", "c")
	require.NotNil(t, ss)
}
