package glist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	l := New[int]()
	for i := 0; i < 5; i++ {
		PushFront(l, i)
	}
	for i := 0; i < 5; i++ {
		PushBack(l, i)
	}
	s1 := make([]int, 0, 10)
	Range(l.Front, func(i int) {
		s1 = append(s1, i)
	})
	s2 := make([]int, 0, 10)
	RangeReverse(l.Back, func(i int) {
		s2 = append(s2, i)
	})
	t.Logf("s1=%v\n", s1)
	t.Logf("s2=%v\n", s2)

	assert := assert.New(t)
	assert.Equal(s2, s1)

	Remove(l, l.Back)
	Remove(l, l.Front)

	lenL := 0
	Range(l.Front, func(i int) {
		lenL++
	})
	assert.Equal(8, lenL)

	assert.Equal(3, l.Front.Value)
	assert.Equal(3, l.Back.Value)
}
