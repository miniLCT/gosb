package gnumerics

import "github.com/miniLCT/gosb/gogenerics/gconstraints"

// Max returns the maximum of two values
func Max[T gconstraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// MaxCollection returns the maximum value in a collection. Returns zero value when collection is empty
func MaxCollection[T gconstraints.Ordered](collection []T) T {
	var maxC T
	l := len(collection)
	if l == 0 {
		return maxC
	}

	maxC = collection[0]
	for i := 1; i < l; i++ {
		c := collection[i]
		if c > maxC {
			maxC = c
		}
	}
	return maxC
}

// MaxBy returns the maximum value in a collection, using the given comparison function
// If several values are equivalent, the first one is returned
// Returns zero value when collection is empty
func MaxBy[T any](collection []T, cmp func(a T, b T) bool) T {
	var maxC T

	l := len(collection)
	if l == 0 {
		return maxC
	}

	maxC = collection[0]
	for i := 1; i < l; i++ {
		c := collection[i]
		if cmp(c, maxC) {
			maxC = c
		}
	}
	return maxC
}

// Min returns the minimum of two values.
func Min[T gconstraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// MinCollection returns the minimum value in a collection. Returns zero value when collection is empty
func MinCollection[T gconstraints.Ordered](collection []T) T {
	var minC T
	l := len(collection)
	if l == 0 {
		return minC
	}

	minC = collection[0]
	for i := 1; i < l; i++ {
		c := collection[i]
		if c < minC {
			minC = c
		}
	}
	return minC
}

// MinBy returns the minimum value in a collection, using the given comparison function
// If several values are equivalent, the first one is returned
// Returns zero value when collection is empty
func MinBy[T any](collection []T, cmp func(a T, b T) bool) T {
	var minC T

	l := len(collection)
	if l == 0 {
		return minC
	}

	minC = collection[0]
	for i := 1; i < l; i++ {
		c := collection[i]
		if cmp(c, minC) {
			minC = c
		}
	}
	return minC
}
