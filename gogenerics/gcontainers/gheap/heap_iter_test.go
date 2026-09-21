package gheap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	less := func(a, b int) bool { return a < b }
	h := New[int](less)
	for _, v := range []int{5, 3, 8, 1} {
		Push(h, v)
	}
	assert.Equal(4, h.Len())
	assert.False(h.Empty())

	top, ok := h.Peek()
	assert.True(ok)
	assert.Equal(1, top)

	var got []int
	for v := range h.All() {
		got = append(got, v)
	}
	assert.Equal([]int{1, 3, 5, 8}, got)
	// All pops every element
	assert.Equal(0, h.Len())
	assert.True(h.Empty())
}

func TestTryPopOnEmptyHeap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := New[string](func(a, b string) bool { return a < b })
	v, ok := TryPop(h)
	assert.False(ok)
	assert.Equal("", v)

	_, ok = h.Peek()
	assert.False(ok)
}

func TestAllBreak(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := New[int](func(a, b int) bool { return a < b })
	for _, v := range []int{5, 3, 8, 1} {
		Push(h, v)
	}

	var got []int
	for v := range h.All() {
		got = append(got, v)
		break
	}
	assert.Equal([]int{1}, got)
	assert.Equal(3, h.Len())
}
