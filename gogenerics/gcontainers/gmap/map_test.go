package gmap

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
)

func TestKeys(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mp1 := map[string]any{
		"a":   1,
		"b":   "b",
		"c":   []string{"c"},
		"100": map[string]any{},
		"11":  true,
	}

	mp2 := map[float64]any{
		1.1:   "1",
		100.0: 114,
		-514:  191.8810,
		-0:    true,
	}

	keys1 := Keys(mp1)
	keys2 := Keys(mp2)
	// TODO: maybe can use this package func to sort
	sort.Strings(keys1)
	stdKeys1 := []string{"100", "11", "a", "b", "c"}
	sort.Float64s(keys2)
	stdKeys2 := []float64{-514, -0, 1.1, 100.00000}

	assert.Equal(stdKeys1, keys1)
	assert.Equal(stdKeys2, keys2)
}

func TestValues(t *testing.T) {
	t.Parallel()
	// assert := assert.New(t)

	mp1 := map[string]any{
		"a":   1,
		"b":   "b",
		"c":   []string{"c"},
		"100": map[string]any{},
		"11":  true,
	}

	mp2 := map[float64]any{
		1.1:   "1",
		100.0: 114,
		-514:  191.8810,
		-0:    true,
	}

	values1 := Values(mp1)
	values2 := Values(mp2)
	// TODO: fix it
	t.Logf("%v\n", values1)
	t.Logf("%v\n", values2)
	// fmt.Println(values1...)
	// fmt.Println(values2)
	// assert.Equal([]any{1, "b", []string{"c"}, map[string]any{}, true}, values1)
	// assert.Equal([]any{"1", 114, 191.881, true}, values2)
}

func TestCopy(t *testing.T) {
	assert := assert.New(t)

	mp1 := map[string]any{
		"a":   1,
		"b":   "b",
		"c":   []string{"c"},
		"100": map[string]any{},
		"11":  true,
	}
	mp2 := Copy(mp1)

	assert.Equal(mp1, mp2)
	assert.NotEqual(fmt.Sprintf("%p", &mp1), fmt.Sprintf("%p", &mp2), "copy should not be the same pointer")
}

func TestMap2Entries(t *testing.T) {
	assert := assert.New(t)

	mp := map[int64]any{
		-1: "-1",
		0:  0,
		1:  true,
		2:  []float64{2.0},
		3:  map[string]any{"3": 3},
	}
	entries := Map2Entries(mp)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	stdEntries := []gconstraints.Entry[int64, any]{
		{Key: -1, Value: "-1"},
		{Key: 0, Value: 0},
		{Key: 1, Value: true},
		{Key: 2, Value: []float64{2.0}},
		{Key: 3, Value: map[string]any{"3": 3}},
	}

	assert.Equal(stdEntries, entries)
}

func TestEntries2Map(t *testing.T) {
	assert := assert.New(t)

	entries := []gconstraints.Entry[int64, any]{
		{Key: -1, Value: "-1"},
		{Key: 0, Value: 0},
		{Key: 1, Value: true},
		{Key: 2, Value: []float64{2.0}},
		{Key: 3, Value: map[string]any{"3": 3}},
	}
	mp := Entries2Map(entries)
	stdMp := map[int64]any{
		-1: "-1",
		0:  0,
		1:  true,
		2:  []float64{2.0},
		3:  map[string]any{"3": 3},
	}

	assert.Equal(stdMp, mp)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	mp1 := map[string]string{
		"a": "1",
		"b": "b",
	}
	mp2 := map[string]string{
		"a": "1",
	}
	mp3 := map[string]string{
		"a": "1",
		"b": "b",
	}

	assert.False(Equal(mp1, mp2))
	assert.False(Equal(mp2, mp3))
	assert.True(Equal(mp1, mp3))
}

func TestLen(t *testing.T) {
	assert := assert.New(t)

	mp := map[string]string{
		"a": "***",
		"b": "b",
	}
	assert.Equal(2, Len(mp))
	mp["c"] = "c"
	assert.Equal(3, Len(mp))
	delete(mp, "ddd")
	assert.Equal(3, Len(mp))
	delete(mp, "c")
	assert.Equal(2, Len(mp))
}

func TestEqualWithFunc(t *testing.T) {
	assert := assert.New(t)

	mp1 := map[string]string{
		"a": "1",
		"b": "b",
	}
	mp2 := map[string]string{
		"a": "1",
		"b": "bbb",
	}
	mp3 := Copy(mp1)

	assert.Equal(false, EqualWithFunc(mp1, mp2, func(a, b string) bool {
		return a == b
	}))
	assert.Equal(true, EqualWithFunc(mp1, mp3, func(a, b string) bool {
		return a == b
	}))
}

func TestContains(t *testing.T) {
	assert := assert.New(t)

	type Node struct {
		Name string
	}
	mp := map[Node]any{
		{Name: "a"}: "1",
	}
	assert.Equal(true, Contains(mp, Node{Name: "a"}))
	assert.Equal(false, Contains(mp, Node{Name: "b"}))
	Clear(mp)
	assert.Equal(false, Contains(mp, Node{Name: "a"}))
}

func TestClear(t *testing.T) {
	assert := assert.New(t)

	mp := map[string]any{
		"a": "1",
		"b": 123,
	}
	t.Logf("now len: %d\n", len(mp))
	t.Logf("mp[%v]===%v", "a", mp["a"])
	t.Logf("mp[%v]===%v", "b", mp["b"])
	t.Logf("mp[%v]===%v", "c", mp["c"])
	Clear(mp)
	t.Logf("after clear len: %d\n", len(mp))
	t.Logf("mp[%v]===%v", "a", mp["a"])
	assert.Equal(0, len(mp))
	assert.Equal(map[string]any{}, mp)
}
