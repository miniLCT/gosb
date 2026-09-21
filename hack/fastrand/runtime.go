package fastrand

import (
	_ "unsafe" // for go:linkname
)

// The functions below are bound to the per-M random source of the Go runtime,
// they are much faster than math/rand because they need no locking and no
// seed injection.

//go:linkname runtimefastrand runtime.fastrand
func runtimefastrand() uint32

//go:linkname runtimefastrand64 runtime.fastrand64
func runtimefastrand64() uint64

//go:linkname runtimefastrandu runtime.fastrandu
func runtimefastrandu() uint
