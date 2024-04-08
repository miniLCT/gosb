package rbt

import (
	"errors"
	"fmt"
	"math"
	"strings"

	gconstraints "github.com/miniLCT/gosb/gogenerics/gconstraints"
)

const (
	black bool = true
	red   bool = false
)

var (
	colorMap = map[bool]string{
		true:  "B",
		false: "R",
	}
)

// rbNode represents the node for a red black tree.

type rbNode[K comparable, V any] struct {
	parent *rbNode[K, V] // Parent of the node.
	left   *rbNode[K, V] // Left child of the node.
	right  *rbNode[K, V] // Right child of the node.
	color  bool          // Color of the node.
	key    K             // Key of the node.
	value  V             // Value of the node.
}

// newRedBlackNode creates and returns a red black node with the specified key and value.
func newRedBlackNode[K comparable, V any](key K, value V, sentinel *rbNode[K, V]) *rbNode[K, V] {
	return &rbNode[K, V]{
		parent: sentinel,
		left:   sentinel,
		right:  sentinel,
		key:    key,
		value:  value,
	}
}

// String returns a string of the form (key, value, color) representing the node.
func (node rbNode[K, V]) String() string {
	return fmt.Sprintf("(%v, %v, %v)", node.key, node.value, colorMap[node.color])
}

// RbTree implementation of a red black tree in which each node has a key and associate value.

type RbTree[K comparable, V any] struct {
	root     *rbNode[K, V]        // The root of the tree.
	sentinel *rbNode[K, V]        // The sentinel node.
	len      int                  // Number of nodes in the tree.
	less     gconstraints.Less[K] // The comparison for ordering keys.
}

// New creates a RbTree. Keys are compared using the less function which should satisfy.
// k1 < k2 => less(k1, k2) = true and less(k2,k1) = false.
// k1 = k2 => less(k1,k2) = false and less(k2,k1) = false.
// k1 > k2 -> less(k1,k2) = false and less(k2,k1) = true.
func New[K comparable, V any](less gconstraints.Less[K]) *RbTree[K, V] {
	sentinel := rbNode[K, V]{parent: nil, left: nil, right: nil, color: black}
	return &RbTree[K, V]{
		root:     &sentinel,
		less:     less,
		sentinel: &sentinel,
	}
}

// Insert inserts a node of the form (key,value) into the tree. If the key already exist its value will be updated,
// the currently stored value is returned.
func (tree *RbTree[K, V]) Insert(key K, value V) V {
	node := newRedBlackNode(key, value, tree.sentinel)
	stored, ok := tree.insert(node)
	if ok {
		tree.insertFix(node)
		tree.len++
		return gconstraints.Empty[V]()
	}
	return stored
}

// Update replaces the value stored in the node with given key and returns the previous value that was stored.
func (tree *RbTree[K, V]) Update(key K, value V) (V, bool) {
	node := tree.search(key)
	if node == tree.sentinel {
		return node.value, false
	}
	temp := node.value
	node.value = value
	return temp, true
}

// insert inserts a node into the tree. For internal use to support Insert function.
func (tree *RbTree[K, V]) insert(z *rbNode[K, V]) (V, bool) {
	var y *rbNode[K, V] = tree.sentinel
	x := tree.root
	for x != tree.sentinel {
		y = x
		if z.key == x.key {
			stored := x.value
			x.value = z.value
			return stored, false
		} else if tree.less(z.key, x.key) {
			x = x.left
		} else {
			x = x.right
		}
	}
	z.parent = y
	if y == tree.sentinel {
		tree.root = z
	} else if tree.less(z.key, y.key) {
		y.left = z
	} else {
		y.right = z
	}
	z.color = red
	return tree.sentinel.value, true
}

// insertFix fixes the tree after an insertion. For internal use to support Insert function.
func (tree *RbTree[K, V]) insertFix(z *rbNode[K, V]) {
	var y *rbNode[K, V]
	for !z.parent.color {
		if z.parent == z.parent.parent.left {
			y = z.parent.parent.right
			if !y.color {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					tree.leftRotate(z)
				}
				z.parent.color = black
				z.parent.parent.color = red
				tree.rightRotate(z.parent.parent)
			}
		} else {
			y = z.parent.parent.left
			if !y.color {
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					tree.rightRotate(z)
				}

				z.parent.color = black
				z.parent.parent.color = red
				tree.leftRotate(z.parent.parent)
			}
		}
	}
	tree.root.color = black
}

// leftRotate performs a left rotation around node x of the tree. For internal use to support deleteFix and insertFix functions.
func (tree *RbTree[K, V]) leftRotate(x *rbNode[K, V]) {
	y := x.right
	x.right = y.left

	if y.left != tree.sentinel {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == tree.sentinel {
		tree.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.left = x
	x.parent = y
}

// rightRotate performs a right rotation around the node x of the tree. For internal use to support deleteFix and insertFix functions.
func (tree *RbTree[K, V]) rightRotate(x *rbNode[K, V]) {
	y := x.left
	x.left = y.right

	if y.right != tree.sentinel {
		y.right.parent = x
	}
	y.parent = x.parent
	if x.parent == tree.sentinel {
		tree.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}
	y.right = x
	x.parent = y
}

// transplant performs transplant operation on the tree. For internal use to support deleteFix and insertFix functions.
func (tree *RbTree[K, V]) transplant(u *rbNode[K, V], v *rbNode[K, V]) {
	if u.parent == tree.sentinel {
		tree.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

// minimum returns the node with the smallest key value in the tree. For internal use to support Minimum and Delete functions.
func (tree *RbTree[K, V]) minimum(node *rbNode[K, V]) *rbNode[K, V] {
	if node.left == tree.sentinel {
		return node
	}
	return tree.minimum(node.left)
}

// search finds the node with the given key in the tree. For internal use to support Search function.
func (tree *RbTree[K, V]) search(key K) *rbNode[K, V] {
	x := tree.root
	for x != tree.sentinel {
		if x.key == key {
			return x
		} else if tree.less(x.key, key) {
			x = x.right
		} else {
			x = x.left
		}
	}
	return x
}

// SubTree returns a new tree that consists of nodes with keys that are in the specified key range [fromKey,toKey]. If fromInclusive is
// true then range includes fromKey otherwise it is left out and if toInclusive is true toKey is included in the range.
func (tree *RbTree[K, V]) SubTree(fromKey K, fromInclusive bool, toKey K, toInclusive bool) *RbTree[K, V] {
	if tree.less(toKey, fromKey) && !(toKey == fromKey) {
		panic(errors.New("undefined range lower key cannot be greater than upper key bound"))
	}
	subTree := New[K, V](tree.less)
	var traverse func(node *rbNode[K, V])
	traverse = func(node *rbNode[K, V]) {
		if node == tree.sentinel {
			return
		}

		if node.left != tree.sentinel {
			traverse(node.left)
		}

		if node.key == fromKey && fromInclusive {
			subTree.Insert(node.key, node.value)
		} else if node.key == toKey && toInclusive {
			subTree.Insert(node.key, node.value)
		} else if tree.less(fromKey, node.key) && tree.less(node.key, toKey) {
			subTree.Insert(node.key, node.value)
		}

		if node.right != tree.sentinel {
			traverse(node.right)
		}
	}
	traverse(tree.root)
	return subTree
}

// LeftSubTree returns a new tree that consists of nodes with keys that are less than or equals the specified key. If inclusive is
// true then the node with an equal key is included otherwise its left out.
func (tree *RbTree[K, V]) LeftSubTree(key K, inclusive bool) *RbTree[K, V] {
	subTree := New[K, V](tree.less)
	var traverse func(node *rbNode[K, V])
	traverse = func(node *rbNode[K, V]) {
		if node == tree.sentinel {
			return
		}

		if node.left != tree.sentinel {
			traverse(node.left)
		}
		if node.key == key && inclusive {
			subTree.Insert(node.key, node.value)
		} else if tree.less(node.key, key) {
			subTree.Insert(node.key, node.value)
		}

		if node.right != tree.sentinel {
			traverse(node.right)
		}
	}
	traverse(tree.root)
	return subTree
}

// RightSubTree returns a new tree that consists of nodes with keys that are greater than or equals than the specified key. If inclusive is
// true then the node with an equal key is included otherwise its left out.
func (tree *RbTree[K, V]) RightSubTree(key K, inclusive bool) *RbTree[K, V] {
	subTree := New[K, V](tree.less)
	var traverse func(node *rbNode[K, V])
	traverse = func(node *rbNode[K, V]) {
		if node == tree.sentinel {
			return
		}

		if node.left != tree.sentinel {
			traverse(node.left)
		}
		if node.key == key && inclusive {
			subTree.Insert(node.key, node.value)
		} else if tree.less(key, node.key) {
			subTree.Insert(node.key, node.value)
		}

		if node.right != tree.sentinel {
			traverse(node.right)
		}
	}
	traverse(tree.root)
	return subTree
}

// Search checks if the tree contains a node with the specified key.
func (tree *RbTree[K, V]) Search(key K) bool {
	return tree.search(key) != tree.sentinel
}

// Get returns the value of the node with the given key.
func (tree *RbTree[K, V]) Get(key K) V {
	node := tree.search(key)
	if node == tree.sentinel {
		return gconstraints.Empty[V]()
	}
	return node.value
}

// GetIf returns the values of the nodes with keys that satisfy the given predicate.
func (tree *RbTree[K, V]) GetIf(f func(K) bool) []V {
	values := make([]V, 0)
	var traverse func(node *rbNode[K, V])
	traverse = func(node *rbNode[K, V]) {
		if node == tree.sentinel {
			return
		}

		if node.left != tree.sentinel {
			traverse(node.left)
		}

		if f(node.key) {
			values = append(values, node.value)
		}

		if node.right != tree.sentinel {
			traverse(node.right)
		}
	}
	traverse(tree.root)
	return values
}

// Delete deletes the node with the specified key from the tree and returns the value that was stored.
func (tree *RbTree[K, V]) Delete(key K) V {
	node := tree.search(key)
	if node == tree.sentinel {
		return gconstraints.Empty[V]()
	}
	tree.delete(node)
	tree.len = int(math.Max(0, float64(tree.len-1)))
	node.left = nil
	node.right = nil
	temp := node.value
	node = nil
	return temp
}

// delete the node z from the tree. For internal use to support Delete function.
func (tree *RbTree[K, V]) delete(z *rbNode[K, V]) {
	var x, y *rbNode[K, V]
	y = z
	yOriginalColor := y.color
	if z.left == tree.sentinel {
		x = z.right
		tree.transplant(z, z.right)
	} else if z.right == tree.sentinel {
		x = z.left
		tree.transplant(z, z.left)
	} else {
		y = tree.minimum(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			x.parent = y
		} else {
			tree.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}

		tree.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	if yOriginalColor {
		tree.deleteFix(x)
	}
}

// deleteFix fixes the tree after a delete operation. For internal use to support Delete function.
func (tree *RbTree[K, V]) deleteFix(x *rbNode[K, V]) {
	var s *rbNode[K, V]
	for x != tree.root && x.color {
		if x == x.parent.left {
			s = x.parent.right
			if !s.color {
				s.color = black
				x.parent.color = red
				tree.leftRotate(x.parent)
				s = x.parent.right
			}

			if s.left.color && s.right.color {
				s.color = red
				x = x.parent
			} else {
				if s.right.color {
					s.left.color = black
					s.color = red
					tree.rightRotate(s)
					s = x.parent.right
				}

				s.color = x.parent.color
				x.parent.color = black
				s.right.color = black
				tree.leftRotate(x.parent)
				x = tree.root
			}
		} else {
			s = x.parent.left
			if !s.color {
				s.color = black
				x.parent.color = red
				tree.rightRotate(x.parent)
				s = x.parent.left
			}

			if s.left.color && s.right.color {
				s.color = red
				x = x.parent
			} else {
				if s.left.color {
					s.right.color = black
					s.color = red
					tree.leftRotate(s)
					s = x.parent.left
				}

				s.color = x.parent.color
				x.parent.color = black
				s.left.color = black
				tree.rightRotate(x.parent)
				x = tree.root
			}
		}
	}
	x.color = black
}

// values collects all the values in the tree into a slice using an in order traversal. For internal use to support Values function.
func (tree *RbTree[K, V]) values(node *rbNode[K, V], data []V, index *int) {
	if node == tree.sentinel {
		return
	}
	if node.left != tree.sentinel {
		tree.values(node.left, data, index)
	}
	data[*index] = node.value
	*index++
	if node.right != tree.sentinel {
		tree.values(node.right, data, index)
	}
}

// Values returns a slice of values stored in the nodes of the tree using an in order traversal.
func (tree *RbTree[K, V]) Values() []V {
	data := make([]V, tree.len)
	index := 0
	tree.values(tree.root, data, &index)
	return data
}

// nodes collects the nodes of the tree using an in order traversal. For internal use to support Nodes function.
func (tree *RbTree[K, V]) nodes(node *rbNode[K, V], nodes []gconstraints.Pair[K, V], index *int) {
	if node == tree.sentinel {
		return
	}
	if node.left != tree.sentinel {
		tree.nodes(node.left, nodes, index)
	}
	nodes[*index] = gconstraints.Pack(node.key, node.value)
	*index++
	if node.right != tree.sentinel {
		tree.nodes(node.right, nodes, index)
	}
}

// Nodes returns the nodes of the tree using an in order traversal.
func (tree *RbTree[K, V]) Nodes() []gconstraints.Pair[K, V] {
	nodes := make([]gconstraints.Pair[K, V], tree.len)
	index := 0
	tree.nodes(tree.root, nodes, &index)
	return nodes
}

// keys collects the keys in the tree into a slice using an in order traversal. For internal use to support Keys function.
func (tree *RbTree[K, V]) keys(node *rbNode[K, V], data []K, index *int) {
	if node == tree.sentinel {
		return
	}
	if node.left != tree.sentinel {
		tree.keys(node.left, data, index)
	}
	data[*index] = node.key
	*index++
	if node.right != tree.sentinel {
		tree.keys(node.right, data, index)
	}
}

// Keys returns a slice of the keys in the tree using an in order traversal.
func (tree *RbTree[K, V]) Keys() []K {
	data := make([]K, tree.len)
	index := 0
	tree.keys(tree.root, data, &index)
	return data
}

// Len returns the size of the tree.
func (tree *RbTree[K, V]) Len() int {
	return tree.len
}

// Clear deletes all the nodes in the tree.
func (tree *RbTree[K, V]) Clear() {
	tree.root = nil
	tree.sentinel = nil
	tree.len = 0
	sentinel := &rbNode[K, V]{parent: nil, left: nil, right: nil, color: black}
	tree.root = sentinel
	tree.sentinel = sentinel
}

// Empty checks if the tree is empty.
func (tree *RbTree[K, V]) Empty() bool {
	return tree.len == 0
}

// printInOrder a helper for string formatting the tree for pretty printing. For internal use to support String function.
func (tree *RbTree[K, V]) printInOrder(node *rbNode[K, V], sb *strings.Builder) {
	if node == tree.sentinel {
		return
	}
	if node.left != tree.sentinel {
		tree.printInOrder(node.left, sb)
	}
	sb.WriteString(fmt.Sprint(node) + " ")
	if node.right != tree.sentinel {
		tree.printInOrder(node.right, sb)
	}
}

// String for pretty printing the tree.
func (tree *RbTree[K, V]) String() string {
	var sb strings.Builder

	tree.printInOrder(tree.root, &sb)
	return "{" + strings.TrimSpace(sb.String()) + "}"
}
