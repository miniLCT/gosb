package gmap

import "github.com/miniLCT/gosb/gogenerics/gconstraints"

// Keys returns a slice of keys from the map. Note that the keys will be an indeterminate order
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of values from the map. Note that the values will be an indeterminate order
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))

	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Copy returns a shallow copy of this map
func Copy[K comparable, V any](m map[K]V) map[K]V {
	res := make(map[K]V, len(m))

	for k, v := range m {
		res[k] = v
	}
	return res
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
	for k := range m {
		delete(m, k)
	}
}

// Map2Entries transforms a map into slice of key-value pairs
func Map2Entries[K comparable, V any](m map[K]V) []gconstraints.Entry[K, V] {
	entries := make([]gconstraints.Entry[K, V], 0, len(m))

	for k, v := range m {
		entries = append(entries, gconstraints.Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return entries
}

// Entries2Map transforms a slice of key-value pairs into a map
func Entries2Map[K comparable, V any](entries []gconstraints.Entry[K, V]) map[K]V {
	m := make(map[K]V, len(entries))

	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}

// Equal returns whether two maps contain the same key-value pairs
func Equal[K, V comparable](m1, m2 map[K]V) bool {
	if Len(m1) != Len(m2) {
		return false
	}

	for k, v1 := range m1 {
		v2, ok := m1[k]
		if !ok || v2 != v1 {
			return false
		}
	}
	return true
}

// EqualWithFunc returns whether two maps contain the same key-value pairs with the given equal function
func EqualWithFunc[K comparable, V1, V2 any](m1 map[K]V1, m2 map[K]V2, eqFunc func(V1, V2) bool) bool {
	if Len(m1) != Len(m2) {
		return false
	}

	for k, v1 := range m1 {
		v2, ok := m2[k]
		if !ok || !eqFunc(v1, v2) {
			return false
		}
	}
	return true
}
