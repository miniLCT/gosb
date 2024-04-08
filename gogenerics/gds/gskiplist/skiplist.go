package gskiplist

import (
	"fmt"
	"github.com/miniLCT/gosb/gogenerics/gconstraints"
	"math/bits"

	"github.com/miniLCT/gosb/hack/fastrand"
)

const (
	// maxLevel denotes the maximum height of the skiplist. This height will keep the skiplist
	// efficient for up to 34m entries. If there is a need for much more, please adjust this constant accordingly.
	maxLevel = 25
)

type SkipListItfc[K gconstraints.Ordered, V any] interface {
	IsEmpty() bool
	Find(key K) (V, bool)
	FindGreaterOrEqual(key K) (*SkipListElement[K, V], bool)
	Delete(key K)
	Insert(key K, value V)
	GetSmallestNode() *SkipListElement[K, V]
	GetLargestNode() *SkipListElement[K, V]
	String() string
	ChangeValue(key K, newValue V) bool
	GetNodeCount() int
}

var _ SkipListItfc[int, int] = (*SkipList[int, int])(nil)

// var _ SkipListItfc[int, string] = (*SkipList[int, string])(nil)

// SkipListElement represents one actual Node in the skiplist structure.
// It saves the actual element, pointers to the next nodes and a pointer to one previous node.

type SkipListElement[K gconstraints.Ordered, V any] struct {
	next  [maxLevel]*SkipListElement[K, V]
	level int
	key   K
	value V
	prev  *SkipListElement[K, V]
}

// SkipList is the actual skiplist representation.
// It saves all nodes accessible from the start and end and keeps track of element count, levels.
type SkipList[K gconstraints.Ordered, V any] struct {
	startLevels  [maxLevel]*SkipListElement[K, V]
	endLevels    [maxLevel]*SkipListElement[K, V]
	maxNewLevel  int
	maxLevel     int
	elementCount int
}

// New returns a new empty, initialized Skiplist.
func New[K gconstraints.Ordered, V any]() *SkipList[K, V] {
	list := &SkipList[K, V]{
		startLevels:  [maxLevel]*SkipListElement[K, V]{},
		endLevels:    [maxLevel]*SkipListElement[K, V]{},
		maxNewLevel:  maxLevel,
		maxLevel:     0,
		elementCount: 0,
	}

	return list
}

// IsEmpty checks, if the skiplist is empty.
func (t *SkipList[K, V]) IsEmpty() bool {
	return t.startLevels[0] == nil
}

func (t *SkipList[K, V]) generateLevel(maxLevel int) int {
	level := maxLevel - 1
	// First we apply some mask which makes sure that we don't get a level
	// above our desired level. Then we find the first set bit.
	x := fastrand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes <= maxLevel {
		level = zeroes
	}
	return level
}

func (t *SkipList[K, V]) findEntryIndex(key K, level int) int {
	// Find good entry point, so we don't accidentally skip half the list.
	for i := t.maxLevel; i >= 0; i-- {
		if t.startLevels[i] != nil && t.startLevels[i].key <= key || i <= level {
			return i
		}
	}
	return 0
}

func (t *SkipList[K, V]) findExtended(key K, findGreaterOrEqual bool) (*SkipListElement[K, V], bool) {
	if t.IsEmpty() {
		return nil, false
	}

	index := t.findEntryIndex(key, 0)
	var currentNode *SkipListElement[K, V]

	currentNode = t.startLevels[index]
	nextNode := currentNode

	// In case, that our first element is already greater-or-equal.
	if findGreaterOrEqual && currentNode.key > key {
		return currentNode, true
	}

	for {
		if currentNode.key == key {
			return currentNode, true
		}

		nextNode = currentNode.next[index]

		// Which direction are we continuing next time?
		if nextNode != nil && nextNode.key <= key {
			// Go right
			currentNode = nextNode
		} else {
			if index > 0 {
				// Early exit
				if currentNode.next[0] != nil && currentNode.next[0].key == key {
					return currentNode.next[0], true
				}
				// Go down
				index--
			} else {
				// Element is not found and we reached the bottom.
				if findGreaterOrEqual {
					return nextNode, nextNode != nil
				}
				return nil, false
			}
		}
	}
}

// Find tries to find an element in the skiplist based on the key from the given ListElement.
// elem can be used, if ok is true.
// Find Expected Complexity O(log(n))
func (t *SkipList[K, V]) Find(key K) (V, bool) {
	if t == nil {
		return gconstraints.Empty[V](), false
	}

	if elem, ok := t.findExtended(key, false); ok {
		return elem.value, true
	}
	return gconstraints.Empty[V](), false
}

// FindGreaterOrEqual finds the first element, that is greater or equal to the given ListElement e.
// The comparison is done on the keys.
// FindGreaterOrEqual Expected Complexity O(log(n))
func (t *SkipList[K, V]) FindGreaterOrEqual(key K) (*SkipListElement[K, V], bool) {
	if t == nil {
		return nil, false
	}

	return t.findExtended(key, true)
}

// Delete removes an element equal to e from the skiplist, if there is one.
// If there are multiple entries with the same value, Delete will remove one of them
// (Which one will change based on the actual skiplist layout)
// Delete Expected Complexity O(log(n))
func (t *SkipList[K, V]) Delete(key K) {
	if t == nil || t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)

	var currentNode *SkipListElement[K, V]
	nextNode := currentNode
	for {
		if currentNode == nil {
			nextNode = t.startLevels[index]
		} else {
			nextNode = currentNode.next[index]
		}

		// Found and remove.
		if nextNode != nil && nextNode.key == key {
			if currentNode != nil {
				currentNode.next[index] = nextNode.next[index]
			}
			if index == 0 {
				if nextNode.next[index] != nil {
					nextNode.next[index].prev = currentNode
				}
				t.elementCount--
			}
			// Link from start needs readjustments.
			if t.startLevels[index] == nextNode {
				t.startLevels[index] = nextNode.next[index]
				// currently the highest node.
				if t.startLevels[index] == nil {
					t.maxLevel = index - 1
				}
			}
			// Link from end needs readjustments.
			if nextNode.next[index] == nil {
				t.endLevels[index] = currentNode
			}
			nextNode.next[index] = nil
		}
		if nextNode != nil && nextNode.key < key {
			// Go right
			currentNode = nextNode
		} else {
			// Go down
			index--
			if index < 0 {
				break
			}
		}
	}
}

// Insert inserts the given ListElement into the skiplist.
// Insert Expected Complexity O(log(n)).
//
// gocyclo:ignore
func (t *SkipList[K, V]) Insert(key K, e V) {
	if t == nil {
		return
	}

	if _, ok := t.findExtended(key, false); ok {
		t.ChangeValue(key, e)
		return
	}

	level := t.generateLevel(t.maxNewLevel)

	// Only grow the height of the skiplist by one at a time.
	if level > t.maxLevel {
		level = t.maxLevel + 1
		t.maxLevel = level
	}

	elem := &SkipListElement[K, V]{
		next:  [maxLevel]*SkipListElement[K, V]{},
		level: level,
		value: e,
		key:   key,
	}

	t.elementCount++

	newFirst, newLast := true, true

	elemKey := key
	if !t.IsEmpty() {
		newFirst = elemKey < t.startLevels[0].key
		newLast = elemKey > t.endLevels[0].key
	}

	normallyInserted := false
	if !newFirst && !newLast {
		normallyInserted = true

		index := t.findEntryIndex(elemKey, level)
		var currentNode *SkipListElement[K, V]
		nextNode := t.startLevels[index]
		for {
			if currentNode == nil {
				nextNode = t.startLevels[index]
			} else {
				nextNode = currentNode.next[index]
			}

			// Connect node to next
			if index <= level && (nextNode == nil || nextNode.key > elemKey) {
				elem.next[index] = nextNode
				if currentNode != nil {
					currentNode.next[index] = elem
				}
				if index == 0 {
					elem.prev = currentNode
					if nextNode != nil {
						nextNode.prev = elem
					}
				}
			}

			if nextNode != nil && nextNode.key <= elemKey {
				// Go right
				currentNode = nextNode
			} else {
				// Go down
				index--
				if index < 0 {
					break
				}
			}
		}
	}

	// Where we have a left-most position that needs to be referenced.
	for i := level; i >= 0; i-- {
		didSomething := false
		if newFirst || normallyInserted {
			if t.startLevels[i] == nil || t.startLevels[i].key > elemKey {
				if i == 0 && t.startLevels[i] != nil {
					t.startLevels[i].prev = elem
				}
				elem.next[i] = t.startLevels[i]
				t.startLevels[i] = elem
			}

			// link the endLevels to this element!
			if elem.next[i] == nil {
				t.endLevels[i] = elem
			}
			didSomething = true
		}

		if newLast {
			// Places the element after the very last element on this level.
			// This is very important, so we are not linking the very first element (newFirst AND newLast) to itself.
			if !newFirst {
				if t.endLevels[i] != nil {
					t.endLevels[i].next[i] = elem
				}
				if i == 0 {
					elem.prev = t.endLevels[i]
				}
				t.endLevels[i] = elem
			}

			// Link the startLevels to this element.
			if t.startLevels[i] == nil || t.startLevels[i].key > elemKey {
				t.startLevels[i] = elem
			}
			didSomething = true
		}
		if !didSomething {
			break
		}
	}
}

// GetValue extracts the ListElement value from a skiplist node.
func (e *SkipListElement[K, V]) GetValue() V {
	if e == nil {
		return gconstraints.Empty[V]()
	}
	return e.value
}

// GetSmallestNode returns the very first/the smallest node in the skiplist.
// GetSmallestNode Complexity O(1)
func (t *SkipList[K, V]) GetSmallestNode() *SkipListElement[K, V] {
	return t.startLevels[0]
}

// GetLargestNode returns the very last/the largest node in the skiplist.
// GetLargestNode Complexity O(1)
func (t *SkipList[K, V]) GetLargestNode() *SkipListElement[K, V] {
	return t.endLevels[0]
}

// Next returns the next element based on the given node.
// Next will loop around to the first node, if you call it on the last.
func (t *SkipList[K, V]) Next(e *SkipListElement[K, V]) *SkipListElement[K, V] {
	if e.next[0] == nil {
		return t.startLevels[0]
	}
	return e.next[0]
}

// Prev returns the previous element based on the given node.
// Prev will loop around to the last node, if you call it on the first.
func (t *SkipList[K, V]) Prev(e *SkipListElement[K, V]) *SkipListElement[K, V] {
	if e.prev == nil {
		return t.endLevels[0]
	}
	return e.prev
}

// GetNodeCount returns the number of nodes currently in the skiplist.
func (t *SkipList[K, V]) GetNodeCount() int {
	return t.elementCount
}

// ChangeValue can be used to change the actual value of a node in the skiplist
// without the need of Deleting and reinserting the node again.
// Be advised, that ChangeValue only works, if the actual key will stay the same.
// ok is an indicator, weather the value is actually changed.
func (t *SkipList[K, V]) ChangeValue(key K, newValue V) bool {
	e, ok := t.findExtended(key, false)
	if !ok || e == nil {
		return false
	}

	e.value = newValue
	return true
}

// String returns a string format of the skiplist. Useful to get a graphical overview and/or debugging.
func (t *SkipList[K, V]) String() string {
	s := ""

	s += " --> "
	for i, l := range t.startLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = fmt.Sprintf("%+v", l.value)
		}
		s += fmt.Sprintf("[%v]", next)

		if i == 0 {
			s += "    "
		}
	}
	s += "\n"

	node := t.startLevels[0]
	for node != nil {
		s += fmt.Sprintf("%v: ", node.value)
		for i := 0; i <= node.level; i++ {
			l := node.next[i]
			next := "---"
			if l != nil {
				next = fmt.Sprintf("%+v", l.value)
			}

			if i == 0 {
				prev := "---"
				if node.prev != nil {
					prev = fmt.Sprintf("%+v", node.prev.value)
				}
				s += fmt.Sprintf("[%v|%v]", prev, next)
			} else {
				s += fmt.Sprintf("[%v]", next)
			}
			if i < node.level {
				s += " -> "
			}
		}
		s += "\n"
		node = node.next[0]
	}

	s += " --> "
	for i, l := range t.endLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = fmt.Sprintf("%+v", l.value)
		}
		s += fmt.Sprintf("[%v]", next)
		if i == 0 {
			s += "    "
		}
	}
	s += "\n"
	return s
}
