package gnumerics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resInt := Max(114, 514)
	resFloat := Max(100.0, 99.9)
	resString := Max("109", "11")

	assert.Equal(514, resInt)
	assert.Equal(100.00, resFloat)
	assert.Equal("11", resString)
}

func TestMin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resInt := Min(114, 514)
	resFloat := Min(100.0, 99.9)
	resString := Min("109", "11")

	assert.Equal(114, resInt)
	assert.Equal(99.9, resFloat)
	assert.Equal("109", resString)
}

func TestMaxCollection(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resEmpty := MaxCollection([]int{})
	resInt1 := MaxCollection([]int{1})
	resInt2 := MaxCollection([]int64{1, 2, 3, 4, 5})
	resString := MaxCollection([]string{"11", "109", "100", "a"})
	resFloat := MaxCollection([]float64{-100.0, 99.9, 99.8})

	assert.Equal(0, resEmpty)
	assert.Equal(1, resInt1)
	assert.Equal(int64(5), resInt2)
	assert.Equal("a", resString)
	assert.Equal(99.9, resFloat)
}

func TestMinCollection(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resEmpty := MinCollection([]string{})
	resInt1 := MinCollection([]int{1})
	resInt2 := MinCollection([]int64{1, 2, 3, 4, 5})
	resString := MinCollection([]string{"11", "109", "100", "a"})
	resFloat := MinCollection([]float64{-100.0, 99.9, 99.8})

	assert.Equal("", resEmpty)
	assert.Equal(1, resInt1)
	assert.Equal(int64(1), resInt2)
	assert.Equal("100", resString)
	assert.Equal(-100.0, resFloat)
}

func TestMaxBy(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	res1 := MaxBy([]string{"hello", "the", "cruel", "woooooorld"}, func(a, b string) bool {
		return len(a) > len(b)
	})
	res2 := MaxBy([]string{"hello", "the", "cruel", "world"}, func(a, b string) bool {
		return len(a) > len(b)
	})
	res3 := MaxBy([]int{1, 92, -9, 123, -456}, func(a, b int) bool {
		return a > b
	})

	assert.Equal("woooooorld", res1)
	assert.Equal("hello", res2)
	assert.Equal(123, res3)
}

func TestMinBy(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	res1 := MinBy([]string{"hello", "the", "cruel", "woooooorld"}, func(a, b string) bool {
		return len(a) < len(b)
	})
	res2 := MinBy([]string{"hello", "thhhe", "cruel", "world"}, func(a, b string) bool {
		return len(a) < len(b)
	})
	res3 := MinBy([]int{1, 92, -9, 123, -456}, func(a, b int) bool {
		return a < b
	})

	assert.Equal("the", res1)
	assert.Equal("hello", res2)
	assert.Equal(-456, res3)
}
