// Package gset implements a generic hash set on top of the built-in map.
package gset

import (
	"iter"
	"maps"
	"slices"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
)

// Set is the hashset datastructure
type Set[T comparable] map[T]struct{}

// NewSet returns an empty set
func NewSet[T comparable]() Set[T] {
	s := make(Set[T])
	return s
}

// NewSetWithSize returns an empty set initialized with specific size
func NewSetWithSize[T comparable](size int) Set[T] {
	s := make(Set[T], size)
	return s
}

// NewSetFrom collects the values of seq into a new set.
func NewSetFrom[T comparable](seq iter.Seq[T]) Set[T] {
	s := NewSet[T]()
	for v := range seq {
		s[v] = struct{}{}
	}
	return s
}

// All returns an iterator over the elements of the set.
// Note that the elements will be yielded in an indeterminate order.
func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s)
}

// Len returns the number of elements of this set
func Len[T comparable](s Set[T]) int {
	return len(s)
}

// Contains returns true if this set contains the specified element
func Contains[T comparable](s Set[T], e T) bool {
	_, ok := s[e]
	return ok
}

// Add adds the specified element to this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func Add[T comparable](s Set[T], e T) bool {
	s[e] = struct{}{}
	return true
}

// Remove removes the specified element from this set
// Always returns true due to the build-in map doesn't indicate caller whether the given element already exists
// Reserves the return type for future extension
func Remove[T comparable](s Set[T], e T) bool {
	delete(s, e)
	return true
}

// Items returns a slice of elements from the set. Note that the elements will be an indeterminate order
func Items[T comparable](s Set[T]) []T {
	items := make([]T, 0, len(s))
	for k := range s.All() {
		items = append(items, k)
	}
	return items
}

// SortedItems returns a slice of elements from the set in ascending order
func SortedItems[T gconstraints.Ordered](s Set[T]) []T {
	return slices.Sorted(s.All())
}

// Clear removes all the elements from this set
func Clear[T comparable](s Set[T]) {
	clear(s)
}
