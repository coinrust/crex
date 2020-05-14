package utils

import (
	"sync"
	"time"
)

type IdGenerate struct {
	idHigh int64
	id     int64
	mu     sync.RWMutex
}

func (g *IdGenerate) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.id += 1
	return g.idHigh + g.id
}

func NewIdGenerate(baseTime time.Time) *IdGenerate {
	now := time.Now()
	base := time.Date(2006, 1, 2, 0, 0, 0, 0, now.Location())
	year, month, day := baseTime.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	d := date.Sub(base).Hours() / 24.0
	idHigh := int64(d) * 10000
	return &IdGenerate{idHigh: idHigh}
}
