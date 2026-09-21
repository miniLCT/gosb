package gskiplist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestSkipList() *SkipList[int, string] {
	list := New[int, string]()
	for k := 1; k <= 10; k++ {
		list.Insert(k, "v")
	}
	return list
}

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	list := newTestSkipList()
	var keys []int
	for k := range list.All() {
		keys = append(keys, k)
	}
	assert.Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, keys)
}

func TestBackward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	list := newTestSkipList()
	var keys []int
	for k := range list.Backward() {
		keys = append(keys, k)
	}
	assert.Equal([]int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}, keys)
}

func TestAllOnEmptySkipList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	list := New[int, string]()
	var n int
	for range list.All() {
		n++
	}
	for range list.Backward() {
		n++
	}
	assert.Equal(0, n)
}
