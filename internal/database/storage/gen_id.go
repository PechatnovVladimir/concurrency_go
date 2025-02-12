package storage

import (
	"math"
	"sync/atomic"
)

type ID struct {
	counter atomic.Int64
}

func NewID(previousID int64) *ID {
	generator := &ID{}
	generator.counter.Store(previousID)
	return generator
}

func (g *ID) Generate() int64 {
	g.counter.CompareAndSwap(math.MaxInt64, 0)
	return g.counter.Add(1)
}
