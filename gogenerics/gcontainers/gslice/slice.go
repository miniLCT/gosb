package gslice

import (
	"github.com/miniLCT/gosb/gogenerics/constraints"
	"github.com/miniLCT/gosb/hack/fastrand"
)

// Copy returns a shallow copy of the given slice
func Copy[T any](s []T) []T {
	// Preserve nil in case it matters
	if s == nil {
		return nil
	}
	return append([]T{}, s...)
}

// Equal returns whether two slices are equal: the same length and all
// elements equal. Note that size=0 and nil are considered equal; floating
// point NaNs are not considered equal
func Equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EqualWithFunc returns whether two slices are equal using a comparison function on each pair of elements
func EqualWithFunc[T1, T2 any](a []T1, b []T2, eq func(T1, T2) bool) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v1 := range a {
		v2 := b[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

// Index returns the index of the first occurrence of target in s, or -1 if not present
func Index[T comparable](s []T, target T) int {
	for i := range s {
		if s[i] == target {
			return i
		}
	}
	return -1
}

// IndexWithFunc returns the first index i satisfying eq(s[i]), or -1 if none do
func IndexWithFunc[T any](s []T, eq func(T) bool) int {
	for i := range s {
		if eq(s[i]) {
			return i
		}
	}
	return -1
}

// Contains returns whether target is present in s.
func Contains[T comparable](s []T, target T) bool {
	return Index(s, target) >= 0
}

// ContainsWithFunc return whether at least one element e of s satisfies eq(e).
func ContainsWithFunc[T any](s []T, eq func(T) bool) bool {
	return IndexWithFunc(s, eq) >= 0
}

// Len returns the length of slice
func Len[T any](s []T) int {
	return len(s)
}

// Unique returns a new slice containing only the unique elements of s, in the order they first appear.
func Unique[T comparable](s []T) ([]T, int) {
	l := Len(s)
	if l == 0 {
		return make([]T, 0), 0
	}

	idx := 0
	seem := make(map[T]struct{}, len(s)) // comparable instead of any
	uniqS := make([]T, len(s))
	for _, v := range s {
		if _, ok := seem[v]; ok {
			continue
		}
		uniqS[idx] = v
		idx++
		seem[v] = struct{}{}
	}
	return uniqS[:idx], idx
}

// UniqueWithFunc returns a new slice containing only the unique elements
// satisfying func f of s, in the order they first appear.
func UniqueWithFunc[T any, U comparable](s []T, f func(T) U) ([]T, int) {
	l := Len(s)
	if l == 0 {
		return make([]T, 0), 0
	}

	idx := 0
	seem := make(map[U]struct{}, len(s))
	uniqS := make([]T, len(s))
	for _, v := range s {
		key := f(v)

		if _, ok := seem[key]; ok {
			continue
		}

		uniqS[idx] = v
		idx++
		seem[key] = struct{}{}
	}
	return uniqS[:idx], idx
}

// Reverse means the first becomes the last, the second becomes the second to last, and so on
func Reverse[T any](s []T) []T {
	l := len(s)
	mid := l >> 1

	for i := 0; i < mid; i++ {
		s[i], s[l-1-i] = s[l-1-i], s[i]
	}
	return s
}

// Merge returns a new slice containing all the elements of s1
// followed by all the elements of s2. If s1 and s2 are both nil,
// returns the empty slice. If either is nil, returns a copy of
// the other.
func Merge[T any](s1, s2 []T) []T {
	if s1 == nil && s2 == nil {
		return []T{}
	}
	if s1 == nil {
		return Copy(s2)
	}
	if s2 == nil {
		return Copy(s1)
	}
	return append(s1, s2...)
}

// IsSorted reports whether x is sorted in ascending order
func IsSorted[T constraints.Ordered](x []T) bool {
	for i := len(x) - 1; i > 0; i-- {
		if x[i] < x[i-1] {
			return false
		}
	}
	return true
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function
func IsSortedFunc[T any](x []T, less constraints.Less[T]) bool {
	for i := len(x) - 1; i > 0; i-- {
		if less(x[i], x[i-1]) {
			return false
		}
	}
	return true
}

// Shuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm
func Shuffle[T any](collection []T) []T {
	fastrand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})
	return collection
}
