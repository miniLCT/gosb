package gmap

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	got := maps.Collect(All(map[string]int{"a": 1, "b": 2}))
	assert.Equal(map[string]int{"a": 1, "b": 2}, got)
}

func TestCollect(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	entries := Collect(func(yield func(string, int) bool) {
		if !yield("a", 1) {
			return
		}
		_ = yield("b", 2)
	})
	assert.Equal(map[string]int{"a": 1, "b": 2}, entries)
}

func TestInsert(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dst := map[string]int{"a": 1}
	Insert(dst, map[string]int{"a": 10, "b": 2})
	assert.Equal(map[string]int{"a": 10, "b": 2}, dst)
}

func TestEqualDetectsDifference(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.True(Equal(map[string]int{"a": 1}, map[string]int{"a": 1}))
	// regression: Equal used to compare the first map with itself
	assert.False(Equal(map[string]int{"a": 1}, map[string]int{"a": 2}))
	assert.False(Equal(map[string]int{"a": 1}, map[string]int{"b": 1}))
	assert.False(Equal(map[string]int{"a": 1}, map[string]int{}))
}

func TestKeysValuesWithIterator(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	m := map[string]int{"a": 1, "b": 2}
	assert.Equal(2, len(Keys(m)))
	assert.ElementsMatch([]int{1, 2}, Values(m))
	assert.ElementsMatch([]string{"a", "b"}, slices.Collect(maps.Keys(m)))
}
