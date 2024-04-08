package gconstraints

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type Node struct{}

	assert.Empty(Empty[string]())
	assert.Empty(Empty[int]())
	assert.Empty(Empty[Node]())
	assert.Empty(Empty[chan string]())
	assert.Empty(Empty[map[string]*Node]())
}

func TestToPtr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type Node struct {
		ID int
	}

	res := ToPtr(
		Node{ID: 100},
	)
	resBool := ToPtr(true)
	trueTmp := true
	resString := ToPtr("hello the cruel world")
	resTmp := "hello the cruel world"

	assert.Equal(&Node{ID: 100}, res)
	assert.NotEqual(res, Node{ID: 99})
	assert.Equal(&trueTmp, resBool)
	assert.Equal(&resTmp, resString)
}

func TestToValue(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	str1 := "abc"
	ptr := &str1

	assert.Equal("abc", ToValue(ptr))
	assert.Equal("", ToValue[string](nil))
	assert.Equal(0, ToValue[int](nil))
	assert.Nil(ToValue[*string](nil))
	assert.Nil(ToValue[map[string]any](nil))
	assert.EqualValues(ptr, ToValue(&ptr))
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type Node struct {
		ID int
	}

	assert.True(IsEmpty(""))
	assert.False(IsEmpty("abc"))
	assert.True(IsEmpty(0))
	assert.False(IsEmpty(1))
	assert.True(IsEmpty(Node{}))
	assert.False(IsEmpty(Node{ID: 123}))
}

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type Node struct {
		ID int
	}

	assert.False(IsNotEmpty(""))
	assert.True(IsNotEmpty("abc"))
	assert.False(IsNotEmpty(0))
	assert.True(IsNotEmpty(1))
	assert.False(IsNotEmpty(Node{}))
	assert.True(IsNotEmpty(Node{ID: 123}))
}
