package gqueue

import (
	"github.com/miniLCT/gosb/gcontainers/glist"
	"github.com/miniLCT/gosb/gogenerics/constraints"
)

// Queue is a simple FIFO queue, not thread-safe

type Queue[T any] struct {
	list   *glist.List[T]
	length int
}

// New returns an empty FIFO queue
func New[T any]() *Queue[T] {
	return &Queue[T]{
		list:   glist.New[T](),
		length: 0,
	}
}

// Len returns the number of items currently in the queue
func Len[T any](q *Queue[T]) int {
	return q.length
}

// Push adds an element to the tail of the queue, reserves the return type for future extension
func Push[T any](q *Queue[T], v T) error {
	glist.PushBack(q.list, v)
	q.length++
	return nil
}

// Pop removes an element from the head of the queue
func Pop[T any](q *Queue[T]) (T, error) {
	if IsEmpty(q) {
		return constraints.Empty[T](), ErrorEmptyQueue
	}
	val := q.list.Front.Value
	glist.Remove(q.list, q.list.Front)
	q.length--
	return val, nil
}

// Peek retrieves but does not remove the head of the queue
func Peek[T any](q *Queue[T]) (T, error) {
	if IsEmpty(q) {
		// todo:return panic or error?
		return constraints.Empty[T](), ErrorEmptyQueue
	}
	return q.list.Front.Value, nil
}

// PeekAll returns all elements in the queue without removing them
func PeekAll[T any](q *Queue[T]) []T {
	res := make([]T, q.length)
	var idx int
	glist.Range(q.list.Front, func(v T) {
		res[idx] = v
		idx++
	})
	return res
}

// IsEmpty returns whether the queue is empty
func IsEmpty[T any](q *Queue[T]) bool {
	return q.length == 0
}

// Clear empties the queue
func Clear[T any](q *Queue[T]) {
	q.list = glist.New[T]()
	q.length = 0
}

// Iterator returns a channel that will be filled with the elements
func Iterator[T any](q *Queue[T]) <-chan T {
	ch := make(chan T, q.length)
	defer close(ch)
	for {
		val, err := Pop(q)
		if err != nil {
			break
		}
		ch <- val
	}
	return ch
}

func Gen[T any, S ~[]T](s S) *Queue[T] {
	q := New[T]()
	for _, v := range s {
		_ = Push(q, v)
	}
	return q
}

// todo: wrap queue with interface
