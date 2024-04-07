package gqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFifo(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	sli := []int{1, 2, 3}
	fifo := Gen(sli)
	assert.False(IsEmpty(fifo))
	assert.Equal(3, Len(fifo))

	e, err := Pop(fifo)
	assert.Equal(1, e)
	assert.Nil(err)
	assert.Equal(2, Len(fifo))

	_ = Push(fifo, 4)
	assert.Equal(3, Len(fifo))

	e, err = Peek(fifo)
	assert.Equal(2, e)
	assert.Nil(err)

	es := PeekAll(fifo)
	assert.Equal([]int{2, 3, 4}, es)
	assert.Nil(err)

	ch := Iterator(fifo)
	ss := make([]int, 0, len(ch))
	assert.Equal(3, len(ch))
	for e := range ch {
		ss = append(ss, e)
	}
	assert.Equal([]int{2, 3, 4}, ss)

	Clear(fifo)
	assert.Equal(0, Len(fifo))
	assert.True(IsEmpty(fifo))

	emptyQueue := New[string]()
	assert.True(IsEmpty(emptyQueue))
	e2, err := Pop(emptyQueue)
	assert.Equal("", e2)
	assert.NotNil(err)
	e2, err = Peek(emptyQueue)
	assert.Equal("", e2)
	assert.NotNil(err)
}
