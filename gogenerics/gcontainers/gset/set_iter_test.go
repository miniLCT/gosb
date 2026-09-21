package gset

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := NewSet[string]()
	Add(s, "hello")
	Add(s, "world")

	var got []string
	for v := range s.All() {
		got = append(got, v)
	}
	assert.ElementsMatch([]string{"hello", "world"}, got)
	assert.Equal(2, Len(s))
}

func TestAllBreak(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := NewSetFrom(slices.Values([]int{1, 2, 3}))
	var seen int
	for range s.All() {
		seen++
		if seen == 2 {
			break
		}
	}
	assert.Equal(2, seen)
	assert.Equal(3, Len(s))
}

func TestNewSetFrom(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := NewSetFrom(slices.Values([]int{1, 2, 2, 3}))
	assert.Equal(3, Len(s))
	assert.True(Contains(s, 1))
	assert.True(Contains(s, 3))
	assert.False(Contains(s, 4))
}

func TestSortedItems(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := NewSetFrom(slices.Values([]int{3, 1, 2}))
	assert.Equal([]int{1, 2, 3}, SortedItems(s))
}

func TestClearBuiltin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := NewSetFrom(slices.Values([]int{1, 2}))
	Clear(s)
	assert.Equal(0, Len(s))
	assert.Empty(Items(s))
}
