package bitset

import (
	"errors"
	"testing"
)

func TestAll(t *testing.T) {
	b := New(0b10100000)

	var set []int
	for offset, truth := range b.All() {
		if truth {
			set = append(set, offset)
		}
	}
	if len(set) != 2 || set[0] != 0 || set[1] != 2 {
		t.Fatalf("All() set bits = %v, want [0 2]", set)
	}
}

func TestAllBreak(t *testing.T) {
	b := New(0xFF, 0xFF)

	var seen int
	for range b.All() {
		seen++
		if seen == 3 {
			break
		}
	}
	if seen != 3 {
		t.Fatalf("All() visited %d bits after break, want 3", seen)
	}
}

func TestSetOffsetOutOfRange(t *testing.T) {
	b := New(0xFF)

	_, err := b.Set(-100, true)
	if !errors.Is(err, ErrOffsetOutOfRange) {
		t.Fatalf("Set(-100) error = %v, want ErrOffsetOutOfRange", err)
	}
}

func TestClearWithBuiltin(t *testing.T) {
	b := New(0xFF)
	b.Clear()

	for _, truth := range b.All() {
		if truth {
			t.Fatal("Clear() left a set bit")
		}
	}
}

func TestBinary(t *testing.T) {
	b := New(0b10100000, 0b00000001)

	// 10100000 00000001
	if got := b.Binary(" "); got != "10100000 00000001" {
		t.Fatalf("Binary() = %q, want %q", got, "10100000 00000001")
	}
	if got := b.Binary(""); got != "1010000000000001" {
		t.Fatalf("Binary() = %q, want %q", got, "1010000000000001")
	}
	if got := (&BitSet{}).Binary(" "); got != "" {
		t.Fatalf("Binary() on an empty bit set = %q, want empty", got)
	}
}
