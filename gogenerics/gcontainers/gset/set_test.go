package gset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	strSet := NewSet[string]()

	for _, v := range []string{"hello", "the", "curel", "world"} {
		Add(strSet, v)
	}

	t.Logf("set contains hello:%v\n", Contains(strSet, "hello"))
	t.Logf("set contains helloo:%v\n", Contains(strSet, "helloo"))

	assert := assert.New(t)
	assert.Equal(4, Len(strSet))

	Remove(strSet, "hello")
	t.Logf("after remove hello, set contains hello:%v\n", Contains(strSet, "hello"))
	assert.Equal(3, Len(strSet))

	items := Items(strSet)
	t.Logf("items:%v\n", items)

	Clear(strSet)
	t.Logf("now set is empty: %v\n", strSet)
	assert.Equal(0, Len(strSet))
}
