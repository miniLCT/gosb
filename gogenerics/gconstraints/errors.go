package gconstraints

import (
	"sync/atomic"
)

// AtomicError defines an atomic error.
type AtomicError struct {
	// err holds a *error, so loading it never needs a type assertion.
	err atomic.Pointer[error]
}

// Set sets the error.
func (ae *AtomicError) Set(err error) {
	if err != nil {
		ae.err.Store(&err)
	}
}

// Load returns the error.
func (ae *AtomicError) Load() error {
	if p := ae.err.Load(); p != nil {
		return *p
	}
	return nil
}

// Swap stores the given error and returns the previous one.
func (ae *AtomicError) Swap(err error) error {
	if err == nil {
		return ae.Load()
	}
	if old := ae.err.Swap(&err); old != nil {
		return *old
	}
	return nil
}
