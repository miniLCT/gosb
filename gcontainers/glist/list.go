package glist

// Node is a node in the linked list
type Node[T any] struct {
	Value      T
	Prev, Next *Node[T]
}

// List is a doubly-linked list
type List[T any] struct {
	Front *Node[T]
	Back  *Node[T]
}

// New returns an empty linked list
func New[T any]() *List[T] {
	return &List[T]{}
}

// PushBack adds v to the end of the list
func PushBack[T any](l *List[T], v T) {
	PushBackNode(l, &Node[T]{
		Value: v,
	})
}

// PushBackNode adds the node nd to the back of the list
func PushBackNode[T any](l *List[T], nd *Node[T]) {
	nd.Next = nil
	nd.Prev = l.Back
	if l.Back != nil {
		l.Back.Next = nd
	} else {
		l.Front = nd
	}
	l.Back = nd
}

// PushFront adds v to the beginning of the list
func PushFront[T any](l *List[T], v T) {
	PushFrontNode(l, &Node[T]{
		Value: v,
	})
}

// PushFrontNode adds the node nd to the beginning of the list
func PushFrontNode[T any](l *List[T], nd *Node[T]) {
	nd.Next = l.Front
	nd.Prev = nil
	if l.Front != nil {
		l.Front.Prev = nd
	} else {
		l.Back = nd
	}
	l.Front = nd
}

// Remove removes the node nd from the list
func Remove[T any](l *List[T], nd *Node[T]) {
	if nd.Next != nil {
		nd.Next.Prev = nd.Prev
	} else {
		l.Back = nd.Prev
	}

	if nd.Prev != nil {
		nd.Prev.Next = nd.Next
	} else {
		l.Front = nd.Next
	}
}

// Range iterates over the list and calls f for each element
func Range[T any](nd *Node[T], f func(T)) {
	for nd != nil {
		f(nd.Value)
		nd = nd.Next
	}
}

// RangeReverse iterates over the list in reverse order and calls f for each element
func RangeReverse[T any](nd *Node[T], f func(T)) {
	for nd != nil {
		f(nd.Value)
		nd = nd.Prev
	}
}
