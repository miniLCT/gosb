package sync

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/miniLCT/gosb/hack/fastrand"
)

func TestConcurrentRange(t *testing.T) {
	const mapSize = 1 << 10

	m := Map[int64, int64]{}
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if fastrand.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	if testing.Short() {
		iters = 16
	}
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(k, v int64) bool {
			if v%k != 0 {
				t.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				t.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true
			return true
		})

		if len(seen) != mapSize {
			t.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}
}

func TestIssue40999(t *testing.T) {
	var m sync.Map

	// Since the miss-counting in missLocked (via Delete)
	// compares the miss count with len(m.dirty),
	// add an initial entry to bias len(m.dirty) above the miss count.
	m.Store(nil, struct{}{})

	var finalized uint32

	// Set finalizers that count for collected keys. A non-zero count
	// indicates that keys have not been leaked.
	for atomic.LoadUint32(&finalized) == 0 {
		p := new(int)
		runtime.SetFinalizer(p, func(*int) {
			atomic.AddUint32(&finalized, 1)
		})
		m.Store(p, struct{}{})
		m.Delete(p)
		runtime.GC()
	}
}

func TestMapRangeNestedCall(t *testing.T) { // Issue 46399
	var m sync.Map
	for i, v := range [3]string{"hello", "world", "Go"} {
		m.Store(i, v)
	}
	m.Range(func(key, value any) bool {
		m.Range(func(key, value any) bool {
			// We should be able to load the key offered in the Range callback,
			// because there are no concurrent Delete involved in this tested map.
			if v, ok := m.Load(key); !ok || !reflect.DeepEqual(v, value) {
				t.Fatalf("Nested Range loads unexpected value, got %+v want %+v", v, value)
			}

			// We didn't keep 42 and a value into the map before, if somehow we loaded
			// a value from such a key, meaning there must be an internal bug regarding
			// nested range in the Map.
			if _, loaded := m.LoadOrStore(42, "dummy"); loaded {
				t.Fatalf("Nested Range loads unexpected value, want store a new value")
			}

			// Try to Store then LoadAndDelete the corresponding value with the key
			// 42 to the Map. In this case, the key 42 and associated value should be
			// removed from the Map. Therefore any future range won't observe key 42
			// as we checked in above.
			val := "sync.Map"
			m.Store(42, val)
			if v, loaded := m.LoadAndDelete(42); !loaded || !reflect.DeepEqual(v, val) {
				t.Fatalf("Nested Range loads unexpected value, got %v, want %v", v, val)
			}
			return true
		})

		// Remove key from Map on-the-fly.
		m.Delete(key)
		return true
	})

	// After a Range of Delete, all keys should be removed and any
	// further Range won't invoke the callback. Hence length remains 0.
	length := 0
	m.Range(func(key, value any) bool {
		length++
		return true
	})

	if length != 0 {
		t.Fatalf("Unexpected sync.Map size, got %v want %v", length, 0)
	}
}

func TestCompareAndSwap_NonExistingKey(t *testing.T) {
	m := &sync.Map{}
	if m.CompareAndSwap(m, nil, 42) {
		// See https://go.dev/issue/51972#issuecomment-1126408637.
		t.Fatalf("CompareAndSwap on an non-existing key succeeded")
	}

	mm := &Map[string, any]{}
	if mm.CompareAndSwap("", nil, "foo") {
		// See https://go.dev/issue/51972#issuecomment-1126408637.
		t.Fatalf("CompareAndSwap on an non-existing key succeeded")
	}
}
