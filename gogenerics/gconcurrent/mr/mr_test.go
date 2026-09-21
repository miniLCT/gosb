package mr

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinish(t *testing.T) {
	var count atomic.Int64
	err := Finish(
		func() error {
			count.Add(1)
			return nil
		},
		func() error {
			count.Add(1)
			return nil
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count.Load())
}

func TestFinishWithError(t *testing.T) {
	target := errors.New("boom")
	err := Finish(func() error {
		return target
	})
	assert.ErrorIs(t, err, target)
}

func TestFinishVoid(t *testing.T) {
	var count atomic.Int64
	FinishVoid(
		func() { count.Add(1) },
		func() { count.Add(1) },
		func() { count.Add(1) },
	)
	assert.Equal(t, int64(3), count.Load())
}

func TestForEach(t *testing.T) {
	var sum atomic.Int64
	ForEach(
		func(source chan<- int) {
			for i := 1; i <= 10; i++ {
				source <- i
			}
		},
		func(item int) {
			sum.Add(int64(item))
		},
	)
	assert.Equal(t, int64(55), sum.Load())
}

func TestMapReduce(t *testing.T) {
	val, err := MapReduce(
		func(source chan<- int) {
			for i := 1; i <= 10; i++ {
				source <- i
			}
		},
		func(item int, writer Writer[int], _ func(error)) {
			writer.Write(item * 2)
		},
		func(pipe <-chan int, writer Writer[int], _ func(error)) {
			var sum int
			for v := range pipe {
				sum += v
			}
			writer.Write(sum)
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, 110, val)
}

func TestMapReduceVoid(t *testing.T) {
	var sum atomic.Int64
	err := MapReduceVoid(
		func(source chan<- int) {
			for i := 1; i <= 4; i++ {
				source <- i
			}
		},
		func(item int, writer Writer[int], _ func(error)) {
			writer.Write(item)
		},
		func(pipe <-chan int, _ func(error)) {
			for v := range pipe {
				sum.Add(int64(v))
			}
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), sum.Load())
}

func TestMapReduceCancel(t *testing.T) {
	target := errors.New("cancel it")
	_, err := MapReduce(
		func(source chan<- int) {
			source <- 1
			source <- 2
		},
		func(item int, _ Writer[int], cancel func(error)) {
			cancel(target)
		},
		func(pipe <-chan int, writer Writer[int], _ func(error)) {},
	)
	assert.ErrorIs(t, err, target)
}

func TestMapReduceChan(t *testing.T) {
	source := make(chan int)
	go func() {
		defer close(source)
		for i := 1; i <= 5; i++ {
			source <- i
		}
	}()

	val, err := MapReduceChan(
		source,
		func(item int, writer Writer[int], _ func(error)) {
			writer.Write(item)
		},
		func(pipe <-chan int, writer Writer[int], _ func(error)) {
			var sum int
			for v := range pipe {
				sum += v
			}
			writer.Write(sum)
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, 15, val)
}
