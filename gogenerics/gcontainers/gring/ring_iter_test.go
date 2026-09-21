package gring

import (
	"testing"
)

func TestAll(t *testing.T) {
	r := New[int](3)
	for i := 1; i <= 3; i++ {
		r.Value = i
		r = r.Next()
	}

	var got []int
	for v := range r.All() {
		got = append(got, v)
	}
	if len(got) != 3 || got[0] != 1 || got[2] != 3 {
		t.Fatalf("All() = %v, want [1 2 3]", got)
	}
}

func TestBackward(t *testing.T) {
	r := New[int](3)
	for i := 1; i <= 3; i++ {
		r.Value = i
		r = r.Next()
	}

	var got []int
	for v := range r.Backward() {
		got = append(got, v)
	}
	if len(got) != 3 || got[0] != 1 || got[1] != 3 || got[2] != 2 {
		t.Fatalf("Backward() = %v, want [1 3 2]", got)
	}
}

func TestAllOnNilRing(t *testing.T) {
	var r *Ring[int]
	var n int
	for range r.All() {
		n++
	}
	if n != 0 {
		t.Fatalf("nil ring yielded %d values, want 0", n)
	}
}

// Do must keep the element type, no more any/type assertion at call site.
func TestDoKeepsType(t *testing.T) {
	r := New[int](2)
	r.Value = 10
	r.Next().Value = 20

	sum := 0
	r.Do(func(v int) {
		sum += v
	})
	if sum != 30 {
		t.Fatalf("Do() sum = %d, want 30", sum)
	}
}
