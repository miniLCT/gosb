package rbt

import (
	"testing"

	"github.com/miniLCT/gosb/gogenerics/gconstraints"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tree := New[int, struct{}](func(i1, i2 int) bool { return i1 < i2 })
	assert.NotNil(t, tree.sentinel)
	assert.True(t, tree.root == tree.sentinel)
	assert.True(t, tree.Empty())
	assert.Equal(t, "{}", tree.String())
}

func TestInsert(t *testing.T) {
	type insertTest struct {
		action       func() *RbTree[int, int]
		expectedTree string
		expectedLen  int
	}

	lessThan := func(i1, i2 int) bool { return i1 < i2 }
	insertTests := []insertTest{
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				return tree
			},
			expectedTree: "{(20, 1, B)}",
			expectedLen:  1,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				return tree
			},
			expectedTree: "{(20, 1, B) (30, 2, R)}",
			expectedLen:  2,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(40, 12)
				return tree
			},
			expectedTree: "{(20, 1, R) (30, 2, B) (40, 12, R)}",
			expectedLen:  3,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(40, 12)
				tree.Insert(10, 3)
				return tree
			},
			expectedTree: "{(10, 3, R) (20, 1, B) (30, 2, B) (40, 12, B)}",
			expectedLen:  4,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(40, 12)
				tree.Insert(10, 3)
				tree.Insert(15, 10)
				return tree
			},
			expectedTree: "{(10, 3, R) (15, 10, B) (20, 1, R) (30, 2, B) (40, 12, B)}",
			expectedLen:  5,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(40, 12)
				tree.Insert(10, 3)
				tree.Insert(15, 10)
				tree.Insert(25, 111)
				return tree
			},
			expectedTree: "{(10, 3, B) (15, 10, R) (20, 1, B) (25, 111, R) (30, 2, B) (40, 12, B)}",
			expectedLen:  6,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(40, 12)
				tree.Insert(10, 3)
				tree.Insert(15, 10)
				tree.Insert(25, 111)
				tree.Insert(24, 11)
				return tree
			},
			expectedTree: "{(10, 3, B) (15, 10, R) (20, 1, R) (24, 11, B) (25, 111, R) (30, 2, B) (40, 12, B)}",
			expectedLen:  7,
		},
		{
			action: func() *RbTree[int, int] {
				tree := New[int, int](lessThan)
				tree.Insert(20, 1)
				tree.Insert(30, 2)
				tree.Insert(30, -2)
				tree.Insert(40, 12)
				tree.Insert(10, 3)
				tree.Insert(15, 10)
				tree.Insert(25, 111)
				tree.Insert(24, 11)
				tree.Insert(24, 9)
				return tree
			},
			expectedTree: "{(10, 3, B) (15, 10, R) (20, 1, R) (24, 9, B) (25, 111, R) (30, -2, B) (40, 12, B)}",
			expectedLen:  7,
		},
	}

	for _, test := range insertTests {
		tree := test.action()
		assert.Equal(t, test.expectedTree, tree.String())
		assert.Equal(t, test.expectedLen, tree.Len())
	}
}

func TestDelete(t *testing.T) {
	lessThan := func(i1, i2 int) bool { return i1 < i2 }
	defaultTree := func() *RbTree[int, int] {
		tree := New[int, int](lessThan)
		tree.Insert(20, 1)
		tree.Insert(30, 1)
		tree.Insert(40, 1)
		tree.Insert(10, 1)
		tree.Insert(15, 1)
		tree.Insert(25, 1)
		tree.Insert(24, 1)
		tree.Insert(21, 1)
		tree.Insert(17, 1)
		tree.Insert(41, 1)
		tree.Insert(39, 1)
		return tree
	}

	type deleteTest struct {
		input        *RbTree[int, int]
		action       func(tree *RbTree[int, int])
		expectedTree string
		expectedLen  int
	}

	deleteTests := []deleteTest{
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				tree.Clear()
				tree.Delete(1)
			},
			expectedTree: "{}",
			expectedLen:  0,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				tree.Delete(10)
			},
			expectedTree: "{(15, 1, B) (17, 1, R) (20, 1, R) (21, 1, B) (24, 1, B) (25, 1, B) (30, 1, R) (39, 1, R) (40, 1, B) (41, 1, R)}",
			expectedLen:  10,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				tree.Delete(10)
				tree.Delete(15)
			},
			expectedTree: "{(17, 1, B) (20, 1, R) (21, 1, B) (24, 1, B) (25, 1, B) (30, 1, R) (39, 1, R) (40, 1, B) (41, 1, R)}",
			expectedLen:  9,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30}
				for _, value := range values {
					tree.Delete(value)
				}
			},
			expectedTree: "{(17, 1, B) (20, 1, R) (21, 1, B) (24, 1, B) (25, 1, B) (39, 1, R) (40, 1, B) (41, 1, R)}",
			expectedLen:  8,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				tree.Delete(10)
				tree.Delete(15)
				tree.Delete(30)
				tree.Delete(24)
			},
			expectedTree: "{(17, 1, B) (20, 1, R) (21, 1, B) (25, 1, B) (39, 1, B) (40, 1, R) (41, 1, B)}",
			expectedLen:  7,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25}
				for _, value := range values {
					tree.Delete(value)
				}
			},
			expectedTree: "{(17, 1, B) (20, 1, R) (21, 1, B) (39, 1, B) (40, 1, B) (41, 1, R)}",
			expectedLen:  6,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25, 39}
				for _, value := range values {
					tree.Delete(value)
				}
			},
			expectedTree: "{(17, 1, B) (20, 1, R) (21, 1, B) (40, 1, B) (41, 1, B)}",
			expectedLen:  5,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25, 39, 41}
				for _, value := range values {
					tree.Delete(value)
				}
			},
			expectedTree: "{(17, 1, B) (20, 1, B) (21, 1, R) (40, 1, B)}",
			expectedLen:  4,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25, 39, 41, 40}
				for _, value := range values {
					tree.Delete(value)
				}
			},
			expectedTree: "{(17, 1, B) (20, 1, B) (21, 1, B)}",
			expectedLen:  3,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25, 39, 41, 40}
				for _, value := range values {
					tree.Delete(value)
				}
				tree.Insert(14, 1)
				tree.Delete(21)
			},
			expectedTree: "{(14, 1, B) (17, 1, B) (20, 1, B)}",
			expectedLen:  3,
		},
		{
			input: defaultTree(),
			action: func(tree *RbTree[int, int]) {
				values := []int{10, 15, 30, 24, 25, 39, 41, 40}
				for _, value := range values {
					tree.Delete(value)
				}
				tree.Insert(14, 1)
				tree.Delete(21)
				tree.Insert(18, 1)
				tree.Insert(23, 1)
				tree.Insert(21, 1)
				tree.Delete(17)
			},
			expectedTree: "{(14, 1, B) (18, 1, B) (20, 1, B) (21, 1, R) (23, 1, B)}",
			expectedLen:  5,
		},
		{
			input: New[int, int](lessThan),
			action: func(tree *RbTree[int, int]) {
				tree.Insert(50, 1)
				tree.Insert(80, 1)
				tree.Insert(90, 1)
				tree.Insert(100, 1)
				tree.Insert(120, 1)
				tree.Insert(140, 1)
				tree.Insert(150, 1)
				tree.Insert(110, 1)
				tree.Insert(122, 1)
				tree.Delete(110)
				tree.Delete(150)
			},
			expectedTree: "{(50, 1, B) (80, 1, R) (90, 1, B) (100, 1, B) (120, 1, B) (122, 1, R) (140, 1, B)}",
			expectedLen:  7,
		},
	}

	for _, test := range deleteTests {
		test.action(test.input)
		assert.Equal(t, test.expectedTree, test.input.String())
		assert.Equal(t, test.expectedLen, test.input.Len())
	}
}

func TestSearch(t *testing.T) {
	type searchTest struct {
		input     *RbTree[int, int]
		searchKey int
		expected  bool
	}

	tree := New[int, int](func(i1, i2 int) bool { return i1 < i2 })
	tree.Insert(20, 1)
	tree.Insert(30, 2)
	tree.Insert(40, 3)
	tree.Insert(10, 4)
	tree.Insert(15, 5)
	tree.Insert(25, 6)
	tree.Insert(24, 7)
	tree.Insert(21, 8)
	tree.Insert(17, 9)
	tree.Insert(41, 10)
	tree.Insert(39, 11)

	searchTests := []searchTest{
		{
			input:     New[int, int](func(i1, i2 int) bool { return i1 < i2 }),
			searchKey: 1,
			expected:  false,
		},
		{
			input:     tree,
			searchKey: 1,
			expected:  false,
		},
		{
			input:     tree,
			searchKey: 21,
			expected:  true,
		},
		{
			input:     tree,
			searchKey: 39,
			expected:  true,
		},
		{
			input:     tree,
			searchKey: 24,
			expected:  true,
		},
	}

	for _, test := range searchTests {
		assert.Equal(t, test.expected, test.input.Search(test.searchKey))
	}
}

func TestUpdate(t *testing.T) {
	type updateTest struct {
		input    *RbTree[int, string]
		key      int
		value    string
		expected string
	}

	lessThan := func(i1, i2 int) bool { return i1 < i2 }

	defaultTree := func() *RbTree[int, string] {
		tree := New[int, string](lessThan)
		tree.Insert(20, "A")
		tree.Insert(30, "B")
		tree.Insert(40, "C")
		tree.Insert(10, "D")
		tree.Insert(15, "E")
		tree.Insert(25, "F")
		tree.Insert(24, "G")
		tree.Insert(21, "H")
		tree.Insert(17, "I")
		tree.Insert(41, "J")
		tree.Insert(39, "K")
		return tree
	}

	updateTests := []updateTest{
		{
			input:    New[int, string](lessThan),
			key:      0,
			value:    "RR",
			expected: gconstraints.Empty[string](),
		},
		{
			input:    defaultTree(),
			key:      20,
			value:    "RR",
			expected: "RR",
		},
		{
			input:    defaultTree(),
			key:      -1,
			value:    "AB",
			expected: gconstraints.Empty[string](),
		},
	}

	for _, test := range updateTests {
		test.input.Update(test.key, test.value)
		assert.Equal(t, test.expected, test.input.Get(test.key))
	}
}

func TestKeys(t *testing.T) {
	var less gconstraints.Less[int]
	less = func(i1, i2 int) bool { return i1 < i2 }

	type keysTest struct {
		input    func() *RbTree[int, struct{}]
		expected []int
	}

	keyTests := []keysTest{
		{
			input: func() *RbTree[int, struct{}] {
				tree := New[int, struct{}](less)
				return tree
			},
			expected: []int{},
		},
		{
			input: func() *RbTree[int, struct{}] {
				tree := New[int, struct{}](less)
				tree.Insert(1, struct{}{})
				tree.Insert(2, struct{}{})
				tree.Insert(3, struct{}{})
				return tree
			},
			expected: []int{1, 2, 3},
		},
		{
			input: func() *RbTree[int, struct{}] {
				tree := New[int, struct{}](less)
				tree.Insert(3, struct{}{})
				tree.Insert(2, struct{}{})
				tree.Insert(1, struct{}{})
				return tree
			},
			expected: []int{1, 2, 3},
		},
		{
			input: func() *RbTree[int, struct{}] {
				tree := New[int, struct{}](less)
				tree.Insert(10, struct{}{})
				tree.Insert(2, struct{}{})
				tree.Insert(12, struct{}{})
				tree.Insert(11, struct{}{})
				tree.Insert(5, struct{}{})
				tree.Insert(5, struct{}{})
				tree.Insert(7, struct{}{})

				return tree
			},
			expected: []int{2, 5, 7, 10, 11, 12},
		},
	}

	for _, test := range keyTests {
		assert.Equal(t, test.expected, test.input().Keys())
	}
}

func TestValues(t *testing.T) {
	var less gconstraints.Less[int]
	less = func(i1, i2 int) bool { return i1 < i2 }

	type valuesTest struct {
		input    func() *RbTree[int, int]
		expected []int
	}

	valuesTests := []valuesTest{
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				return tree
			},
			expected: []int{},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(1, 1)
				tree.Insert(2, 2)
				tree.Insert(3, 3)
				return tree
			},
			expected: []int{1, 2, 3},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(3, 3)
				tree.Insert(2, 2)
				tree.Insert(1, 1)
				return tree
			},
			expected: []int{1, 2, 3},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(10, 10)
				tree.Insert(2, 2)
				tree.Insert(12, 12)
				tree.Insert(11, 11)
				tree.Insert(5, 5)
				tree.Insert(5, 90)
				tree.Insert(7, 7)

				return tree
			},
			expected: []int{2, 90, 7, 10, 11, 12},
		},
	}

	for _, test := range valuesTests {
		assert.Equal(t, test.expected, test.input().Values())
	}
}

func TestNodes(t *testing.T) {
	var less gconstraints.Less[int]
	less = func(i1, i2 int) bool {
		return i1 < i2
	}

	type valuesTest struct {
		input    func() *RbTree[int, int]
		expected []gconstraints.Pair[int, int]
	}

	valuesTests := []valuesTest{
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				return tree
			},
			expected: []gconstraints.Pair[int, int]{},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(1, 1)
				tree.Insert(2, 2)
				tree.Insert(3, 3)
				return tree
			},
			expected: []gconstraints.Pair[int, int]{
				gconstraints.Pack(1, 1),
				gconstraints.Pack(2, 2),
				gconstraints.Pack(3, 3),
			},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(3, 3)
				tree.Insert(2, 2)
				tree.Insert(1, 1)
				tree.Insert(5, 22)
				return tree
			},
			expected: []gconstraints.Pair[int, int]{
				gconstraints.Pack(1, 1),
				gconstraints.Pack(2, 2),
				gconstraints.Pack(3, 3),
				gconstraints.Pack(5, 22),
			},
		},
	}

	for _, test := range valuesTests {
		assert.Equal(t, test.expected, test.input().Nodes())
	}
}

func TestGetIf(t *testing.T) {
	var less gconstraints.Less[int]
	less = func(i1, i2 int) bool { return i1 < i2 }

	type valuesTest struct {
		input    func() *RbTree[int, int]
		expected []int
	}

	valuesTests := []valuesTest{
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				return tree
			},
			expected: []int{},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(1, 1)
				tree.Insert(2, 2)
				tree.Insert(3, 3)
				tree.Insert(4, 4)
				return tree
			},
			expected: []int{2, 4},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(3, 3)
				tree.Insert(2, 2)
				tree.Insert(1, 1)
				tree.Insert(9, 9)
				return tree
			},
			expected: []int{2},
		},
		{
			input: func() *RbTree[int, int] {
				tree := New[int, int](less)
				tree.Insert(10, 10)
				tree.Insert(2, 2)
				tree.Insert(12, 12)
				tree.Insert(11, 11)
				tree.Insert(5, 5)
				tree.Insert(7, 7)

				return tree
			},
			expected: []int{2, 10, 12},
		},
	}

	f := func(key int) bool {
		return key%2 == 0
	}

	for _, test := range valuesTests[3:] {
		assert.Equal(t, test.expected, test.input().GetIf(f))
	}
}

func TestLeftSubtTree(t *testing.T) {
	lessThan := func(i1, i2 int) bool { return i1 < i2 }

	tree := New[int, string](lessThan)
	tree.Insert(20, "A")
	tree.Insert(30, "B")
	tree.Insert(40, "C")
	tree.Insert(10, "D")
	tree.Insert(15, "E")
	tree.Insert(25, "F")
	tree.Insert(24, "G")
	tree.Insert(21, "H")
	tree.Insert(17, "I")
	tree.Insert(41, "J")
	tree.Insert(39, "K")

	type leftSubTreeTest struct {
		input          *RbTree[int, string]
		key            int
		inclusive      bool
		expectedKeys   []int
		expectedValues []string
	}

	leftSubTreeTests := []leftSubTreeTest{
		{
			input:          New[int, string](lessThan),
			key:            0,
			inclusive:      true,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            0,
			inclusive:      true,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            25,
			inclusive:      false,
			expectedKeys:   []int{10, 15, 17, 20, 21, 24},
			expectedValues: []string{"D", "E", "I", "A", "H", "G"},
		},
		{
			input:          tree,
			key:            25,
			inclusive:      true,
			expectedKeys:   []int{10, 15, 17, 20, 21, 24, 25},
			expectedValues: []string{"D", "E", "I", "A", "H", "G", "F"},
		},
		{
			input:          tree,
			key:            10,
			inclusive:      true,
			expectedKeys:   []int{10},
			expectedValues: []string{"D"},
		},
		{
			input:          tree,
			key:            10,
			inclusive:      false,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            41,
			inclusive:      true,
			expectedKeys:   []int{10, 15, 17, 20, 21, 24, 25, 30, 39, 40, 41},
			expectedValues: []string{"D", "E", "I", "A", "H", "G", "F", "B", "K", "C", "J"},
		},
		{
			input:          tree,
			key:            41,
			inclusive:      false,
			expectedKeys:   []int{10, 15, 17, 20, 21, 24, 25, 30, 39, 40},
			expectedValues: []string{"D", "E", "I", "A", "H", "G", "F", "B", "K", "C"},
		},
	}

	for _, test := range leftSubTreeTests {
		leftSubTree := test.input.LeftSubTree(test.key, test.inclusive)
		assert.Equal(t, test.expectedKeys, leftSubTree.Keys())
		assert.Equal(t, test.expectedValues, leftSubTree.Values())
	}
}

func TestRightSubtTree(t *testing.T) {
	lessThan := func(i1, i2 int) bool { return i1 < i2 }

	tree := New[int, string](lessThan)
	tree.Insert(20, "A")
	tree.Insert(30, "B")
	tree.Insert(40, "C")
	tree.Insert(10, "D")
	tree.Insert(15, "E")
	tree.Insert(25, "F")
	tree.Insert(24, "G")
	tree.Insert(21, "H")
	tree.Insert(17, "I")
	tree.Insert(41, "J")
	tree.Insert(39, "K")

	type rightSubTreeTest struct {
		input          *RbTree[int, string]
		key            int
		inclusive      bool
		expectedKeys   []int
		expectedValues []string
	}

	rightSubTreeTests := []rightSubTreeTest{
		{
			input:          New[int, string](lessThan),
			key:            50,
			inclusive:      true,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            50,
			inclusive:      true,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            25,
			inclusive:      false,
			expectedKeys:   []int{30, 39, 40, 41},
			expectedValues: []string{"B", "K", "C", "J"},
		},
		{
			input:          tree,
			key:            25,
			inclusive:      true,
			expectedKeys:   []int{25, 30, 39, 40, 41},
			expectedValues: []string{"F", "B", "K", "C", "J"},
		},
		{
			input:          tree,
			key:            41,
			inclusive:      true,
			expectedKeys:   []int{41},
			expectedValues: []string{"J"},
		},
		{
			input:          tree,
			key:            41,
			inclusive:      false,
			expectedKeys:   []int{},
			expectedValues: []string{},
		},
		{
			input:          tree,
			key:            10,
			inclusive:      true,
			expectedKeys:   []int{10, 15, 17, 20, 21, 24, 25, 30, 39, 40, 41},
			expectedValues: []string{"D", "E", "I", "A", "H", "G", "F", "B", "K", "C", "J"},
		},
		{
			input:          tree,
			key:            10,
			inclusive:      false,
			expectedKeys:   []int{15, 17, 20, 21, 24, 25, 30, 39, 40, 41},
			expectedValues: []string{"E", "I", "A", "H", "G", "F", "B", "K", "C", "J"},
		},
	}

	for _, test := range rightSubTreeTests {
		rightSubTree := test.input.RightSubTree(test.key, test.inclusive)
		assert.Equal(t, test.expectedKeys, rightSubTree.Keys())
		assert.Equal(t, test.expectedValues, rightSubTree.Values())
	}
}

func TestSubTree(t *testing.T) {
	lessThan := func(i1, i2 int) bool { return i1 < i2 }
	tree := New[int, int](lessThan)
	tree.Insert(1, 1)
	tree.Insert(2, 2)
	tree.Insert(3, 3)
	tree.Insert(4, 4)
	tree.Insert(5, 5)
	tree.Insert(6, 6)
	tree.Insert(7, 7)
	tree.Insert(8, 8)
	tree.Insert(9, 9)
	tree.Insert(10, 10)

	type subTreeTest = struct {
		input         *RbTree[int, int]
		fromKey       int
		fromInclusive bool
		toKey         int
		toInclusive   bool
		expectedKeys  []int
	}

	subTreeTests := []subTreeTest{
		{
			input:        New[int, int](lessThan),
			fromKey:      0,
			toKey:        0,
			expectedKeys: []int{},
		},
		{
			input:         tree,
			fromKey:       3,
			fromInclusive: false,
			toKey:         7,
			toInclusive:   false,
			expectedKeys:  []int{4, 5, 6},
		},
		{
			input:         tree,
			fromKey:       3,
			fromInclusive: true,
			toKey:         7,
			toInclusive:   false,
			expectedKeys:  []int{3, 4, 5, 6},
		},
		{
			input:         tree,
			fromKey:       3,
			fromInclusive: true,
			toKey:         7,
			toInclusive:   true,
			expectedKeys:  []int{3, 4, 5, 6, 7},
		},
		{
			input:         tree,
			fromKey:       3,
			fromInclusive: false,
			toKey:         7,
			toInclusive:   true,
			expectedKeys:  []int{4, 5, 6, 7},
		},
	}

	for _, test := range subTreeTests {
		subTree := test.input.SubTree(test.fromKey, test.fromInclusive, test.toKey, test.toInclusive)
		assert.Equal(t, test.expectedKeys, subTree.Keys())
	}
}
