package glist

import (
	"testing"
)

func TestAll(t *testing.T) {
	l := New[int]()
	for i := 1; i <= 3; i++ {
		l.PushBack(i)
	}

	var got []int
	for e := range l.All() {
		got = append(got, e.Value)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 2 || got[2] != 3 {
		t.Fatalf("All() = %v, want [1 2 3]", got)
	}
}

func TestBackward(t *testing.T) {
	l := New[int]()
	for i := 1; i <= 3; i++ {
		l.PushBack(i)
	}

	var got []int
	for e := range l.Backward() {
		got = append(got, e.Value)
	}
	if len(got) != 3 || got[0] != 3 || got[2] != 1 {
		t.Fatalf("Backward() = %v, want [3 2 1]", got)
	}
}

func TestValues(t *testing.T) {
	l := NewFrom(collectValues([]int{1, 2}))
	if l.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", l.Len())
	}

	var got []int
	for v := range l.Values() {
		got = append(got, v)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("Values() = %v, want [1 2]", got)
	}
}

func TestTransform(t *testing.T) {
	l := NewFrom(collectValues([]int{1, 2, 3}))
	s := l.Transform(func(v int) string {
		if v == 1 {
			return "one"
		}
		if v == 2 {
			return "two"
		}
		return "three"
	})

	var got []string
	for v := range s.Values() {
		got = append(got, v)
	}
	if len(got) != 3 || got[0] != "one" || got[2] != "three" {
		t.Fatalf("Transform() = %v, want [one two three]", got)
	}
}

func TestRemoveReturnsValue(t *testing.T) {
	l := New[int]()
	e := l.PushBack(42)
	// Remove returns the typed value, not any
	var v int = l.Remove(e)
	if v != 42 {
		t.Fatalf("Remove() = %d, want 42", v)
	}
}

func collectValues(s []int) func(func(int) bool) {
	return func(yield func(int) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
