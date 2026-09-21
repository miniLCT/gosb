// Package gnumerics provides numeric helpers. Implementations delegate to the
// modern standard library: the built-in min/max functions (go1.21) and the
// slices package (go1.21), so they keep the same behaviour with less code.
package gnumerics

import (
	"slices"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
)

// Max returns the maximum of two values
func Max[T gconstraints.Ordered](a, b T) T {
	return max(a, b)
}

// MaxCollection returns the maximum value in a collection. Returns zero value when collection is empty
func MaxCollection[T gconstraints.Ordered](collection []T) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.Max(collection)
}

// MaxBy returns the maximum value in a collection, using the given comparison function
// If several values are equivalent, the first one is returned
// Returns zero value when collection is empty
func MaxBy[T any](collection []T, cmp func(a T, b T) bool) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.MaxFunc(collection, greaterToCmp(cmp))
}

// MaxFunc returns the maximum value in a collection, using a three-way
// comparison function as expected by slices.MaxFunc.
// If several values are equivalent, the first one is returned.
// Returns zero value when collection is empty
func MaxFunc[T any](collection []T, cmp func(a T, b T) int) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.MaxFunc(collection, cmp)
}

// Min returns the minimum of two values.
func Min[T gconstraints.Ordered](a, b T) T {
	return min(a, b)
}

// MinCollection returns the minimum value in a collection. Returns zero value when collection is empty
func MinCollection[T gconstraints.Ordered](collection []T) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.Min(collection)
}

// MinBy returns the minimum value in a collection, using the given comparison function
// If several values are equivalent, the first one is returned
// Returns zero value when collection is empty
func MinBy[T any](collection []T, cmp func(a T, b T) bool) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.MinFunc(collection, lessToCmp(cmp))
}

// MinFunc returns the minimum value in a collection, using a three-way
// comparison function as expected by slices.MinFunc.
// If several values are equivalent, the first one is returned.
// Returns zero value when collection is empty
func MinFunc[T any](collection []T, cmp func(a T, b T) int) T {
	if len(collection) == 0 {
		return gconstraints.Empty[T]()
	}
	return slices.MinFunc(collection, cmp)
}

// Clamp returns v restricted to the inclusive range [lo, hi].
// If lo > hi, hi is returned.
func Clamp[T gconstraints.Ordered](v, lo, hi T) T {
	return min(max(v, lo), hi)
}

// greaterToCmp adapts a "a is greater than b" predicate to the three-way
// comparison function used by slices.MaxFunc/slices.MinFunc.
func greaterToCmp[T any](greater func(a, b T) bool) func(a, b T) int {
	return func(a, b T) int {
		switch {
		case greater(a, b):
			return 1
		case greater(b, a):
			return -1
		default:
			return 0
		}
	}
}

// lessToCmp adapts a "a is less than b" predicate to the three-way comparison
// function used by slices.MaxFunc/slices.MinFunc.
func lessToCmp[T any](less func(a, b T) bool) func(a, b T) int {
	return func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	}
}
