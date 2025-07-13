package verbose

import (
	"fmt"
	"testing"
	"time"
)

func TestVerboseLogger(t *testing.T) {
	st := time.Now().Add(10 * time.Second)
	inMs := time.Until(st).Milliseconds()
	fmt.Println(inMs)
}
