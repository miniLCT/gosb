package gslice

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var got []int
	for v := range Values([]int{1, 2, 3}) {
		got = append(got, v)
	}
	assert.Equal([]int{1, 2, 3}, got)
}

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var indexes []int
	var values []string
	for i, v := range All([]string{"a", "b"}) {
		indexes = append(indexes, i)
		values = append(values, v)
	}
	assert.Equal([]int{0, 1}, indexes)
	assert.Equal([]string{"a", "b"}, values)
}

func TestBackward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var indexes []int
	for i := range Backward([]int{1, 2, 3}) {
		indexes = append(indexes, i)
	}
	assert.Equal([]int{2, 1, 0}, indexes)
}

func TestCollect(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]int{1, 2, 3}, Collect(Values([]int{1, 2, 3})))
}

func TestSorted(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	src := []int{3, 1, 2}
	assert.Equal([]int{1, 2, 3}, Sorted(src))
	// the source is left untouched
	assert.Equal([]int{3, 1, 2}, src)
}

func TestSortedFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]string{"ccc", "bb", "a"}, SortedFunc(
		[]string{"a", "bb", "ccc"},
		func(a, b string) int { return len(b) - len(a) },
	))
}

func TestFilter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	even := Filter([]int{1, 2, 3, 4}, func(v int) bool { return v%2 == 0 })
	assert.Equal([]int{2, 4}, even)
	assert.Equal([]int{}, Filter([]int{}, func(v int) bool { return true }))
}

func TestIterBreak(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// breaking out of the loop must stop the iteration (and the callback)
	var seen int
	for range Values([]int{1, 2, 3, 4}) {
		seen++
		if seen == 2 {
			break
		}
	}
	assert.Equal(2, seen)
	assert.True(slices.Equal([]int{1, 2, 3}, Collect(Values([]int{1, 2, 3}))))
}
