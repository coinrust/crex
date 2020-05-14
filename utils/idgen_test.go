package utils

import (
	"testing"
	"time"
)

func TestIdGenerate(t *testing.T) {
	idGen := NewIdGenerate(time.Now())
	for i := 0; i < 10; i++ {
		id := idGen.Next()
		t.Logf("id=%v", id)
	}
}
