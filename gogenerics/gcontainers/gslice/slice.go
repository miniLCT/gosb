// Package gslice provides helpers for slices. Everything that has an equivalent
// in the standard library (slices, iter) delegates to it, the rest
// (Unique, Shuffle, ...) is implemented here.
package gslice

import (
	"iter"
	"slices"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
	"github.com/miniLCT/gosb/hack/fastrand"
)

// Copy returns a shallow copy of the given slice
func Copy[T any](s []T) []T {
	// slices.Clone preserves nil, just like the documented behaviour here.
	return slices.Clone(s)
}

// Equal returns whether two slices are equal: the same length and all
// elements equal. Note that size=0 and nil are considered equal; floating
// point NaNs are not considered equal
func Equal[T comparable](a, b []T) bool {
	return slices.Equal(a, b)
}

// EqualWithFunc returns whether two slices are equal using a comparison function on each pair of elements
func EqualWithFunc[T1, T2 any](a []T1, b []T2, eq func(T1, T2) bool) bool {
	return slices.EqualFunc(a, b, eq)
}

// Index returns the index of the first occurrence of target in s, or -1 if not present
func Index[T comparable](s []T, target T) int {
	return slices.Index(s, target)
}

// IndexWithFunc returns the first index i satisfying eq(s[i]), or -1 if none do
func IndexWithFunc[T any](s []T, eq func(T) bool) int {
	return slices.IndexFunc(s, eq)
}

// Contains returns whether target is present in s.
func Contains[T comparable](s []T, target T) bool {
	return slices.Contains(s, target)
}

// ContainsWithFunc return whether at least one element e of s satisfies eq(e).
func ContainsWithFunc[T any](s []T, eq func(T) bool) bool {
	return slices.ContainsFunc(s, eq)
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

// Reverse means the first becomes the last, the second becomes the second to last, and so on.
// The reversal happens in place and s is returned for convenience.
func Reverse[T any](s []T) []T {
	slices.Reverse(s)
	return s
}

// Merge returns a new slice containing all the elements of s1
// followed by all the elements of s2. If s1 and s2 are both nil,
// returns the empty slice. If either is nil, returns a copy of
// the other.
func Merge[T any](s1, s2 []T) []T {
	merged := slices.Concat(s1, s2)
	if merged == nil {
		return []T{}
	}
	return merged
}

// IsSorted reports whether x is sorted in ascending order
func IsSorted[T gconstraints.Ordered](x []T) bool {
	return slices.IsSorted(x)
}

// IsSortedFunc reports whether x is sorted in ascending order, with less as the
// comparison function
func IsSortedFunc[T any](x []T, less gconstraints.Less[T]) bool {
	return slices.IsSortedFunc(x, lessToCmp(less))
}

// Shuffle returns an array of shuffled values. Uses the Fisher-Yates shuffle algorithm
func Shuffle[T any](collection []T) []T {
	fastrand.Shuffle(len(collection), func(i, j int) {
		collection[i], collection[j] = collection[j], collection[i]
	})
	return collection
}

// Values returns an iterator over the elements of s.
func Values[T any](s []T) iter.Seq[T] {
	return slices.Values(s)
}

// All returns an iterator over index-value pairs of s.
func All[T any](s []T) iter.Seq2[int, T] {
	return slices.All(s)
}

// Backward returns an iterator over index-value pairs of s, in reverse order.
func Backward[T any](s []T) iter.Seq2[int, T] {
	return slices.Backward(s)
}

// Collect collects the values of seq into a new slice.
func Collect[T any](seq iter.Seq[T]) []T {
	return slices.Collect(seq)
}

// Sorted returns a copy of s sorted in ascending order, s itself is left untouched.
func Sorted[T gconstraints.Ordered](s []T) []T {
	return slices.Sorted(slices.Values(s))
}

// SortedFunc returns a copy of s sorted with the given three-way comparison
// function, s itself is left untouched.
func SortedFunc[T any](s []T, cmp func(a, b T) int) []T {
	return slices.SortedFunc(slices.Values(s), cmp)
}

// Filter returns a new slice holding only the elements of s satisfying keep.
func Filter[T any](s []T, keep func(T) bool) []T {
	if len(s) == 0 {
		return make([]T, 0)
	}
	res := make([]T, 0, len(s))
	for _, v := range s {
		if keep(v) {
			res = append(res, v)
		}
	}
	return res
}

// lessToCmp adapts a "a is less than b" predicate to the three-way comparison
// function used by the slices package.
func lessToCmp[T any](less gconstraints.Less[T]) func(a, b T) int {
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
