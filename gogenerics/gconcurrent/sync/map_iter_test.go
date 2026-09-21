package sync

import (
	"testing"
)

func TestMapAll(t *testing.T) {
	m := Map[int, int]{}
	for i := 0; i < 10; i++ {
		m.Store(i, i*2)
	}

	count := 0
	for k, v := range m.All() {
		if v != k*2 {
			t.Fatalf("Map.All() yields %d => %d, want %d", k, v, k*2)
		}
		count++
	}
	if count != 10 {
		t.Fatalf("Map.All() visited %d entries, want 10", count)
	}
}

func TestMapAllBreak(t *testing.T) {
	m := Map[int, int]{}
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}

	count := 0
	for range m.All() {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Fatalf("Map.All() visited %d entries after break, want 3", count)
	}
}

func TestMapClear(t *testing.T) {
	m := Map[int, int]{}
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}

	m.Clear()

	count := 0
	for range m.All() {
		count++
	}
	if count != 0 {
		t.Fatalf("Map.Clear() left %d entries, want 0", count)
	}
	if _, ok := m.Load(1); ok {
		t.Fatal("Map.Clear() left the entry for key 1")
	}

	// the map is still usable afterwards
	m.Store(42, 42)
	if v, ok := m.Load(42); !ok || v != 42 {
		t.Fatalf("Map.Load(42) = %d, %v after Clear", v, ok)
	}
}
