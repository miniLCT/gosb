// Package gmap provides helpers for maps. Everything that has an equivalent in
// the standard library (maps, slices, iter) delegates to it.
package gmap

import (
	"iter"
	"maps"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
)

// Keys returns a slice of keys from the map. Note that the keys will be an indeterminate order
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range maps.Keys(m) {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values from the map. Note that the values will be an indeterminate order
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for v := range maps.Values(m) {
		values = append(values, v)
	}
	return values
}

// Copy returns a shallow copy of this map
func Copy[K comparable, V any](m map[K]V) map[K]V {
	return maps.Clone(m)
}

// Len returns the number of elements of this map
func Len[K comparable, V any](m map[K]V) int {
	return len(m)
}

// Contains returns whether this map contains the specified key.
// Note that if the key is a pointer, the result may be unexpected
func Contains[K comparable, V any](m map[K]V, e K) bool {
	_, ok := m[e]
	return ok
}

// Clear removes all the elements from this map
func Clear[K comparable, V any](m map[K]V) {
	clear(m)
}

// Map2Entries transforms a map into slice of key-value pairs
func Map2Entries[K comparable, V any](m map[K]V) []gconstraints.Entry[K, V] {
	entries := make([]gconstraints.Entry[K, V], 0, len(m))
	for k, v := range All(m) {
		entries = append(entries, gconstraints.Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return entries
}

// Entries2Map transforms a slice of key-value pairs into a map
func Entries2Map[K comparable, V any](entries []gconstraints.Entry[K, V]) map[K]V {
	return Collect(func(yield func(K, V) bool) {
		for _, e := range entries {
			if !yield(e.Key, e.Value) {
				return
			}
		}
	})
}

// Equal returns whether two maps contain the same key-value pairs
func Equal[K, V comparable](m1, m2 map[K]V) bool {
	return maps.Equal(m1, m2)
}

// EqualWithFunc returns whether two maps contain the same key-value pairs with the given equal function
func EqualWithFunc[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2, eqFunc func(V1, V2) bool) bool {
	return maps.EqualFunc(m1, m2, eqFunc)
}

// All returns an iterator over the key-value pairs of m in an indeterminate order.
func All[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return maps.All(m)
}

// Collect collects key-value pairs from seq into a new map.
func Collect[K comparable, V any](seq iter.Seq2[K, V]) map[K]V {
	return maps.Collect(seq)
}

// Insert copies all key-value pairs of src into dst, overwriting existing keys.
func Insert[K comparable, V any](dst, src map[K]V) {
	maps.Insert(dst, All(src))
}
