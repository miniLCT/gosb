package gslice

import (
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Node struct {
	Val int
}

func TestCopy(t *testing.T) {
	t.Parallel()
	t.Helper()
	assert := assert.New(t)

	nilFunc := func() []int64 {
		return nil
	}
	assert.Nil(Copy(nilFunc()))
	assert.Equal([]string{}, Copy([]string{}))
	t.Logf("%p, %p", []string{""}, Copy([]string{""}))
	assert.Equal([]float64{-1}, Copy([]float64{-1}))
	t.Logf("%p, %p", []float64{-1}, Copy([]float64{-1}))
	assert.Equal([]*Node{{1}, {2}}, Copy([]*Node{{1}, {2}}))
	t.Logf("%p, %p", []*Node{{1}, {2}}, Copy([]*Node{{1}, {2}}))
}

func TestEqual(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	nilFunc := func() []int64 {
		return nil
	}
	assert.True(Equal([]int32{1, 2}, []int32{1, 2}))
	assert.False(Equal([]int32{1, 2}, []int32{1, 3}))
	assert.True(Equal([]int64{}, nilFunc())) // nil and empty are equal
	assert.False(Equal([]string{"a", "a", "a"}, []string{"a"}))
	assert.False(Equal([]float64{math.NaN()}, []float64{math.NaN()}))
}

func TestEqualWithFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	lenEqFunc := func(a, b string) bool {
		return len(a) == len(b)
	}
	lesFunc := func(a, b int) bool {
		return a < b
	}

	assert.True(EqualWithFunc([]string{"abc", "", "dd"}, []string{"zza", "", "pp"}, lenEqFunc))
	assert.False(EqualWithFunc([]string{"abc"}, []string{"zzadd"}, lenEqFunc))

	assert.True(EqualWithFunc([]int{1, 2, 3}, []int{4, 5, 6}, lesFunc))
	assert.False(EqualWithFunc([]int{1, 2, 3}, []int{1, 2, 3}, lesFunc))
}

func TestIndex(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal(0, Index([]int32{1, 2}, 1))
	assert.Equal(-1, Index([]int32{100, 200}, 1))
	assert.Equal(-1, Index([]string{}, ""))
}

func TestIndexWithFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	existFunc := func(a string) bool {
		return len(a) > 0
	}
	more100Func := func(a int) bool {
		return a > 100
	}

	assert.Equal(2, IndexWithFunc([]string{"", "", "dd"}, existFunc))
	assert.Equal(-1, IndexWithFunc([]string{"", "", ""}, existFunc))
	assert.Equal(1, IndexWithFunc([]int{1, 200, 3}, more100Func))
	assert.Equal(-1, IndexWithFunc([]int{1, 2, 3}, more100Func))
}

func TestContains(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.True(Contains([]int32{1, 2}, 1))
	assert.False(Contains([]int32{100, 200}, 1))
	assert.False(Contains([]string{}, ""))
}

func TestContainsWithFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	existFunc := func(a string) bool {
		return len(a) > 0
	}
	more100Func := func(a int) bool {
		return a > 100
	}

	assert.True(ContainsWithFunc([]string{"", "", "dd"}, existFunc))
	assert.False(ContainsWithFunc([]string{"", "", ""}, existFunc))
	assert.True(ContainsWithFunc([]int{1, 200, 3}, more100Func))
	assert.False(ContainsWithFunc([]int{1, 2, 3}, more100Func))
}

func TestLen(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	nilFunc := func() []int64 {
		return nil
	}
	assert.Equal(0, Len(nilFunc()))
	assert.Equal(0, Len([]string{}))
	assert.Equal(1, Len([]float64{-1}))
	assert.Equal(2, Len([]*Node{{1}, {2}}))
}

func TestUnique(t *testing.T) {
	t.Parallel()
	t.Helper()
	assert := assert.New(t)

	nilFunc := func() []int64 {
		return nil
	}
	type utWrap[T comparable] struct {
		Slice []T
		Count int
	}

	res1, cnt1 := Unique(nilFunc())
	assert.Equal(res1, []int64{})
	assert.Equal(cnt1, 0)
	res2, cnt2 := Unique([]string{})
	assert.Equal(res2, []string{})
	assert.Equal(cnt2, 0)
	res3, cnt3 := Unique([]float64{-1})
	assert.Equal(res3, []float64{-1})
	assert.Equal(cnt3, 1)
	res4, cnt4 := Unique([]string{"a", "a", "b", "b", "a"})
	assert.Equal(res4, []string{"a", "b"})
	assert.Equal(cnt4, 2)

	// TODO:why????
	// assert.Equal(
	//	utWrap[string]{Slice: []string{}, Count: 0},
	//	utWrapFunc([]string{}),
	// )
	// assert.Equal(
	//	utWrap[float64]{Slice: []float64{-1}, Count: 1},
	//	utWrapFunc([]float64{-1}),
	// )
	// assert.Equal(
	//	utWrap[string]{Slice: []string{"a", "b"}, Count: 2},
	//	utWrapFunc([]string{"a", "a", "b", "b", "a"}),
	// )
}

// so ugly 🤮 TODO: find a better way to do this

type utWrap[T comparable] struct {
	Slice []T
	Count int
}

// so ugly 🤮 TODO: find a better way to do this
func utWrapFunc[T comparable](s []T) utWrap[T] {
	res, cnt := Unique(s)
	return utWrap[T]{Slice: res, Count: cnt}
}

func TestUniqueWithFunc(t *testing.T) {
	t.Parallel()
	t.Helper()
	assert := assert.New(t)

	s1, count1 := UniqueWithFunc([]int{0, 1, 2, 3, 4, 5, 6}, func(i int) int {
		return i % 3
	})
	s2, count2 := UniqueWithFunc([]int{0, 1, 2, 3, 4, 5, 6, 1000}, func(i int) int {
		if i > 10 {
			return i
		}
		return -1
	})
	assert.Equal([]int{0, 1, 2}, s1)
	assert.Equal(3, count1)
	assert.Equal([]int{0, 1000}, s2)
	assert.Equal(2, count2)
}

func TestReverse(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]int{3, 2, 1}, Reverse([]int{1, 2, 3}))
	assert.Equal([]string{"a", "b", "c", "d"}, Reverse([]string{"d", "c", "b", "a"}))
	assert.Equal([]int64{}, Reverse([]int64{}))
}

func TestMerge(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]int{1, 2, 3, 4, 5, 6}, Merge([]int{1, 2, 3}, []int{4, 5, 6}))
	assert.Equal([]string{"a", "b", "c", "d"}, Merge([]string{"a", "b"}, []string{"c", "d"}))
	assert.Equal([]int64{}, Merge([]int64{}, []int64{}))

	nilf := func() []int64 {
		return nil
	}
	assert.Equal([]int64{}, Merge([]int64{}, nilf()))
	assert.Equal([]int64{}, Merge(nilf(), []int64{}))
	assert.Equal([]int64{}, Merge(nilf(), nilf()))
}

func TestIsSorted(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
	var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3} // nolint: lll

	assert.False(IsSorted(ints[:]))
	assert.False(IsSorted(float64s[:]))

	sort.Slice(ints[:], func(i, j int) bool {
		return ints[i] < ints[j]
	})
	sort.Slice(float64s[:], func(i, j int) bool {
		return float64s[i] < float64s[j]
	})
	t.Logf("after sort: %v", ints)
	t.Logf("after sort: %v", float64s)
	assert.True(IsSorted(ints[:]))
	assert.True(IsSorted(float64s[:]))
}

func TestIsSortedWithFunc(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	var ints = []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}

	assert.False(IsSortedFunc(ints, func(a, b int) bool {
		return a < b
	}))
	assert.False(IsSortedFunc(ints, func(a, b int) bool {
		return a > b
	}))

	sort.Slice(ints, func(i, j int) bool {
		return ints[i] > ints[j]
	})
	assert.False(IsSortedFunc(ints, func(a, b int) bool {
		return a < b
	}))
	assert.True(IsSortedFunc(ints, func(a, b int) bool {
		return a > b
	}))
}

func TestShuffle(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	result2 := Shuffle([]int{})

	is.NotEqual(result1, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	is.Equal(result2, []int{})
}
