package gskiplist

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type bench struct {
	setup func(*testing.B, mapInterface)
	perG  func(b *testing.B, pb *testing.PB, i int, m mapInterface)
}

func benchMap(b *testing.B, bench bench) {
	for _, m := range [...]mapInterface{&DeepCopyMap{}, &RWMutexMap{}, &SyncMap[int, int]{}, New[int, int]()} {
		b.Run(fmt.Sprintf("%T", m), func(b *testing.B) {
			if strings.Contains(b.Name(), "BenchmarkLoadOrStoreUnique/*gskiplist.SkipList[int,int]") {
				b.Skip("BenchmarkLoadOrStoreUnique/*gskiplist.SkipList[int,int] white benchmark")
				return
			}
			m = reflect.New(reflect.TypeOf(m).Elem()).Interface().(mapInterface)
			if bench.setup != nil {
				bench.setup(b, m)
			}

			b.ResetTimer()

			var i int64
			b.RunParallel(func(pb *testing.PB) {
				id := int(atomic.AddInt64(&i, 1) - 1)
				bench.perG(b, pb, id*b.N, m)
			})
		})
	}
}

func BenchmarkLoadMostlyHits(b *testing.B) {
	const hits, misses = 1023, 1

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.Insert(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Find(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Find(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadMostlyMisses(b *testing.B) {
	const hits, misses = 1, 1023

	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			for i := 0; i < hits; i++ {
				m.Insert(i, i)
			}
			// Prime the map to get it into a steady state.
			for i := 0; i < hits*2; i++ {
				m.Find(i % hits)
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Find(i % (hits + misses))
			}
		},
	})
}

func BenchmarkLoadOrStoreUnique(b *testing.B) {
	benchMap(b, bench{
		setup: func(b *testing.B, m mapInterface) {
			if _, ok := m.(*DeepCopyMap); ok {
				b.Skip("DeepCopyMap has quadratic running time.")
			}
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Insert(i, i)
			}
		},
	})
}

func BenchmarkLoadOrStoreCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.Insert(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Insert(0, 0)
			}
		},
	})
}

// BenchmarkAdversarialAlloc tests performance when we store a new value
// immediately whenever the map is promoted to clean and otherwise load a
// unique, missing key.
//
// This forces the Load calls to always acquire the map's mutex.
func BenchmarkAdversarialAlloc(b *testing.B) {
	benchMap(b, bench{
		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			var stores, loadsSinceStore int
			for ; pb.Next(); i++ {
				m.Find(i)
				if loadsSinceStore++; loadsSinceStore > stores {
					m.Insert(i, stores)
					loadsSinceStore = 0
					stores++
				}
			}
		},
	})
}

func BenchmarkDeleteCollision(b *testing.B) {
	benchMap(b, bench{
		setup: func(_ *testing.B, m mapInterface) {
			m.Insert(0, 0)
		},

		perG: func(b *testing.B, pb *testing.PB, i int, m mapInterface) {
			for ; pb.Next(); i++ {
				m.Delete(0)
			}
		},
	})
}

// below is the code contains reference map implementations for beach.

// mapInterface is the interface Map implements.

type mapInterface interface {
	Find(int) (int, bool)
	Insert(key, value int)
	Delete(int)
}

var (
	_ mapInterface = &RWMutexMap{}
	_ mapInterface = &DeepCopyMap{}
)

// RWMutexMap is an implementation of mapInterface using a sync.RWMutex.

type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[int]int
}

func (m *RWMutexMap) Find(key int) (value int, ok bool) {
	m.mu.RLock()
	value, ok = m.dirty[key]
	m.mu.RUnlock()
	return
}

func (m *RWMutexMap) Insert(key, value int) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[int]int)
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) Delete(key int) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

// DeepCopyMap is an implementation of mapInterface using a Mutex and
// atomic.Value.  It makes deep copies of the map on every write to avoid
// acquiring the Mutex in Load.

type DeepCopyMap struct {
	mu    sync.Mutex
	clean atomic.Value
}

func (m *DeepCopyMap) Find(key int) (value int, ok bool) {
	clean, _ := m.clean.Load().(map[int]int)
	value, ok = clean[key]
	return value, ok
}

func (m *DeepCopyMap) Insert(key, value int) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) Delete(key int) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *DeepCopyMap) dirty() map[int]int {
	clean, _ := m.clean.Load().(map[int]int)
	dirty := make(map[int]int, len(clean)+1)
	for k, v := range clean {
		dirty[k] = v
	}
	return dirty
}

// sync.Map

type SyncMap[K comparable, V any] struct {
	dirty sync.Map
}

func (m *SyncMap[K, V]) Find(key K) (value V, ok bool) {
	v, ok := m.dirty.Load(key)
	if !ok {
		var v V
		return v, false
	}

	return v.(V), true
}

func (m *SyncMap[K, V]) Insert(key, value int) {
	m.dirty.Store(key, value)
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.dirty.Delete(key)
}
