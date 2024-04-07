package constraints

// Empty returns an empty value of the given type
func Empty[T any]() T {
	var empty T
	return empty
}

// ToPtr converts a value to a pointer.
func ToPtr[T any](v T) *T {
	return &v
}

// ToValue converts a pointer to a value.
func ToValue[T any](v *T) T {
	if v == nil {
		return Empty[T]()
	}

	return *v
}

// IsEmpty returns true if argument is a empty value
func IsEmpty[T comparable](v T) bool {
	var empty T
	return v == empty
}

// IsNotEmpty returns true if argument is not a empty value
func IsNotEmpty[T comparable](v T) bool {
	var empty T
	return v != empty
}

// Entry defines a key/value pairs.

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Pair represent key, value pair.
type Pair[K any, V any] struct {
	key   K
	value V
}

// Pack create a pair with the given key and value
func Pack[K any, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{key: key, value: value}
}

// Unpack returns values contained in tuple.
func (t Pair[K, V]) Unpack() (K, V) {
	return t.key, t.value
}
