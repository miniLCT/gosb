package gpool

import (
	"sync"
)

// A Pool is a generic wrapper around a sync.Pool.
type Pool[T any] struct {
	pool sync.Pool
}

// New creates a new Pool with the provided new function.
//
// The equivalent sync.Pool construct is "sync.Pool{New: fn}"
func New[T any](fn func() T) Pool[T] {
	return Pool[T]{
		pool: sync.Pool{New: func() interface{} { return fn() }},
	}
}

// Get is a generic wrapper around sync.Pool's Get method.
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put is a generic wrapper around sync.Pool's Put method.
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}

type ArrayPool[T any] struct {
	pool sync.Pool
}

// NewArrayPool creates a new arrayPool with the provided new Function.
//
// The equivalent sync.Pool construct is "sync.Pool{New: fn}"
func NewArrayPool[T any](size uint) ArrayPool[T] {
	return ArrayPool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]T, 0, size)
			},
		},
	}
}

// Get is a generic wrapper around sync.Pool's Get method.
func (p *ArrayPool[T]) Get() []T {
	return p.pool.Get().([]T)
}

// Put is a generic wrapper around sync.Pool's Put method.
func (p *ArrayPool[T]) Put(x []T) {
	p.pool.Put(x)
}
