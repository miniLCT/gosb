package gqueue

import (
	"iter"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
	"github.com/miniLCT/gosb/gogenerics/gcontainers/glist"
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
func (q *Queue[T]) Len() int {
	return q.length
}

// Push adds an element to the tail of the queue, reserves the return type for future extension
func (q *Queue[T]) Push(v T) error {
	q.list.PushBack(v)
	q.length++
	return nil
}

// Pop removes an element from the head of the queue
func (q *Queue[T]) Pop() (T, error) {
	if IsEmpty(q) {
		return gconstraints.Empty[T](), ErrorEmptyQueue
	}
	val := q.list.Front().Value
	q.list.Remove(q.list.Front())
	q.length--
	return val, nil
}

// Peek retrieves but does not remove the head of the queue
func Peek[T any](q *Queue[T]) (T, error) {
	if IsEmpty(q) {
		// todo:return panic or error?
		return gconstraints.Empty[T](), ErrorEmptyQueue
	}
	return q.list.Front().Value, nil
}

// PeekAll returns all elements in the queue without removing them
func PeekAll[T any](q *Queue[T]) []T {
	res := make([]T, 0, q.length)
	for v := range q.All() {
		res = append(res, v)
	}
	return res
}

// All returns an iterator over the elements of the queue, from head to tail.
// The queue is left untouched.
func (q *Queue[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for e := q.list.Front(); e != nil; e = e.Next() {
			if !yield(e.Value) {
				return
			}
		}
	}
}

// Drain returns an iterator over the elements of the queue, from head to tail.
// Every yielded element is popped from the queue.
func (q *Queue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			v, err := q.Pop()
			if err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
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

// Iterator returns a channel that will be filled with the elements.
// The elements are popped from the queue while the channel is filled.
//
// Deprecated: use the All (non-destructive) or Drain (destructive) iterators,
// they are allocation free and composable with the slices package.
func Iterator[T any](q *Queue[T]) <-chan T {
	ch := make(chan T, q.length)
	for v := range q.Drain() {
		ch <- v
	}
	close(ch)
	return ch
}

func Gen[T any, S ~[]T](s S) *Queue[T] {
	q := New[T]()
	for _, v := range s {
		_ = q.Push(v)
	}
	return q
}
