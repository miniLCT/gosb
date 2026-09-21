package rbt

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestTree() *RbTree[int, string] {
	tree := New[int, string](func(a, b int) bool { return a < b })
	for _, k := range []int{5, 3, 8, 1, 4} {
		tree.Insert(k, strconv.Itoa(k))
	}
	return tree
}

func TestAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var keys []int
	var values []string
	for k, v := range newTestTree().All() {
		keys = append(keys, k)
		values = append(values, v)
	}
	assert.Equal([]int{1, 3, 4, 5, 8}, keys)
	assert.Equal([]string{"1", "3", "4", "5", "8"}, values)
}

func TestBackward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var keys []int
	for k := range newTestTree().Backward() {
		keys = append(keys, k)
	}
	assert.Equal([]int{8, 5, 4, 3, 1}, keys)
}

func TestAllMatchesKeys(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tree := newTestTree()
	var keys []int
	for k := range tree.All() {
		keys = append(keys, k)
	}
	assert.Equal(tree.Keys(), keys)
}

func TestMinMax(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tree := newTestTree()
	k, v, ok := tree.Min()
	assert.True(ok)
	assert.Equal(1, k)
	assert.Equal("1", v)

	k, v, ok = tree.Max()
	assert.True(ok)
	assert.Equal(8, k)
	assert.Equal("8", v)

	empty := New[int, string](func(a, b int) bool { return a < b })
	_, _, ok = empty.Min()
	assert.False(ok)
	_, _, ok = empty.Max()
	assert.False(ok)
}
