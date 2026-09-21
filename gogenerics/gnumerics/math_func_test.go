package gnumerics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	byLen := func(a, b string) int { return len(a) - len(b) }
	assert.Equal("woooooorld", MaxFunc([]string{"hello", "the", "woooooorld"}, byLen))
	assert.Equal("hello", MaxFunc([]string{"hello", "cruel"}, byLen))
	// zero value on empty input
	assert.Equal("", MaxFunc([]string{}, byLen))
	assert.Equal(0, MaxFunc([]int{}, func(a, b int) int { return a - b }))
}

func TestMinFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	byLen := func(a, b string) int { return len(a) - len(b) }
	assert.Equal("the", MinFunc([]string{"hello", "the", "woooooorld"}, byLen))
	// equal length, the first one wins
	assert.Equal("hello", MinFunc([]string{"hello", "cruel"}, byLen))
	// zero value on empty input
	assert.Equal("", MinFunc([]string{}, byLen))
	assert.Equal(0, MinFunc([]int{}, func(a, b int) int { return a - b }))
}

func TestClamp(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal(0, Clamp(-1, 0, 5))
	assert.Equal(3, Clamp(3, 0, 5))
	assert.Equal(5, Clamp(10, 0, 5))
	assert.Equal(1.5, Clamp(1.5, 0.0, 2.0))
}
