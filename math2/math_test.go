package math2

import (
	"testing"
)

func TestR(t *testing.T) {
	x := ToFixedE5(10000.51)
	if x != 10000.5 {
		t.Error("error")
	}
}
