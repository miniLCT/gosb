package fastrand

import (
	"testing"
)

func TestIntn(t *testing.T) {
	const n = 100
	for i := 0; i < 1000; i++ {
		if v := Intn(n); v < 0 || v >= n {
			t.Fatalf("Intn(%d) = %d, out of range", n, v)
		}
	}
}

func TestUint32n(t *testing.T) {
	const n = uint32(97)
	for i := 0; i < 1000; i++ {
		if v := Uint32n(n); v >= n {
			t.Fatalf("Uint32n(%d) = %d, out of range", n, v)
		}
	}
}

func TestPerm(t *testing.T) {
	const n = 64
	p := Perm(n)
	if len(p) != n {
		t.Fatalf("Perm(%d) returned %d elements", n, len(p))
	}
	seen := make([]bool, n)
	for _, v := range p {
		if v < 0 || v >= n {
			t.Fatalf("Perm(%d) contains %d, out of range", n, v)
		}
		if seen[v] {
			t.Fatalf("Perm(%d) contains %d twice", n, v)
		}
		seen[v] = true
	}
}

func TestShuffleSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ShuffleSlice(s)
	if len(s) != 10 {
		t.Fatalf("ShuffleSlice changed the length to %d", len(s))
	}
	seen := make(map[int]struct{}, len(s))
	for _, v := range s {
		seen[v] = struct{}{}
	}
	if len(seen) != len(s) {
		t.Fatalf("ShuffleSlice lost elements: %v", s)
	}
}

func TestRead(t *testing.T) {
	for _, size := range []int{0, 1, 7, 8, 9, 32} {
		p := make([]byte, size)
		n, err := Read(p)
		if err != nil {
			t.Fatalf("Read(%d) returned error %v", size, err)
		}
		if n != size {
			t.Fatalf("Read(%d) = %d bytes", size, n)
		}
		for i, b := range p {
			if b != 0 {
				// extremely unlikely with a good generator, but keep it simple
				break
			}
			_ = i
		}
	}
}

func TestWyrandUint32UsesReceiver(t *testing.T) {
	r := wyrand(1)
	// Uint32 must advance the state of the receiver, so two calls differ
	if a, b := r.Uint32(), r.Uint32(); a == b {
		t.Fatalf("wyrand.Uint32() returned %d twice", a)
	}
}
