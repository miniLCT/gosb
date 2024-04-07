package gheap

import "github.com/miniLCT/gosb/gogenerics/constraints"

// Heap is the generics implementation of heap

type Heap[T any] struct {
	data []T
	less constraints.Less[T]
}

// New constructs a new heap
func New[T any](less constraints.Less[T]) *Heap[T] {
	return &Heap[T]{
		data: make([]T, 0),
		less: less,
	}
}

// NewWithData build a heap tree with data
func NewWithData[T any](data []T, less constraints.Less[T]) *Heap[T] {
	h := New(less)
	h.data = data
	heapSort(h.data, less)
	return h
}

// siftDown implements the heap property on v[lo:hi].
func siftDown[T any](x []T, index int, less constraints.Less[T]) {
	for {
		left := (index * 2) + 1
		right := left + 1
		if left >= len(x) {
			break
		}
		c := left
		if len(x) > right && less(x[right], x[left]) {
			c = right
		}
		if less(x[index], x[c]) {
			break
		}
		x[c], x[index] = x[index], x[c]
		index = c
	}
}

func siftUp[T any](x []T, index int, less constraints.Less[T]) {
	for index > 0 {
		p := (index - 1) / 2
		if less(x[p], x[index]) {
			break
		}
		x[p], x[index] = x[index], x[p]
		index = p
	}
}

// heapSort is min-heap sort
func heapSort[T any](v []T, less constraints.Less[T]) {
	n := len(v)
	for i := n/2 - 1; i >= 0; i-- {
		siftDown(v, i, less)
	}

	// Build heap with greatest element at top.
	// for i := (len(v) - 1) / 2; i >= 0; i-- {
	// 	siftDown(v, i, len(v), less)
	// }

	// // Pop elements into end of v.
	// for i := len(v) - 1; i >= 1; i-- {
	// 	v[0], v[i] = v[i], v[0]
	// 	siftDown(v[:i], 0, len(v), less) // BUG
	// }
}

// Push pushes the element v onto the heap.
func Push[T any](h *Heap[T], v T) {
	x := &h.data
	(*x) = append((*x), v)
	siftUp(*x, len(*x)-1, h.less)
	h.data = (*x)
}

// Pop removes the minimum element from the heap and returns it.
func Pop[T any](h *Heap[T]) T {
	x := &h.data
	ret := (*x)[0]
	(*x)[0], *x = (*x)[len(*x)-1], (*x)[:len(*x)-1]
	if len(*x) > 0 {
		siftDown((*x), 0, h.less)
	}
	h.data = (*x)
	return ret
}
