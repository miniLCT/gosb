package gqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := Gen([]int{1, 2, 3})
	var got []int
	for v := range q.All() {
		got = append(got, v)
	}
	assert.Equal([]int{1, 2, 3}, got)
	// All must not consume the queue
	assert.Equal(3, q.Len())
}

func TestPeekAllFromIterator(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]int{1, 2, 3}, PeekAll(Gen([]int{1, 2, 3})))
}

func TestDrain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := Gen([]int{1, 2, 3})
	var got []int
	for v := range q.Drain() {
		got = append(got, v)
	}
	assert.Equal([]int{1, 2, 3}, got)
	assert.Equal(0, q.Len())
	assert.True(IsEmpty(q))
}

// Iterator used to close the channel before filling it, which made the send
// panic as soon as the buffer was full.
func TestIteratorDoesNotPanic(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := Gen([]int{1, 2, 3})
	ch := Iterator(q)

	var got []int
	for v := range ch {
		got = append(got, v)
	}
	assert.Equal([]int{1, 2, 3}, got)
	assert.Equal(0, q.Len())
}
