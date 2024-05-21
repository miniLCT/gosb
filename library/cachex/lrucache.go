package cachex

import (
	"container/list"
	"context"
	"sync"
	"time"
)

// NewLRUCacheV2 creates a new LRU cache object.
//
// Parameters:
//
//	caption: maximum capacity, should be > 0, panic if 0
//	maxUsed: maximum usage count, should be >=0, if 0, no usage count limitation
func NewLRUCacheV2[K comparable, V any](caption int, maxUsed int) *LRUCacheV2[K, V] {
	if caption <= 0 {
		panic("Size of LRUCacheV2 should not less than zero")
	}
	return &LRUCacheV2[K, V]{
		lruList: list.New(),
		lruMap:  make(map[any]*item2[V], caption),
		maxUsed: maxUsed,
		caption: caption,
	}
}

var _ Cache[string, any] = (*LRUCacheV2[string, any])(nil)

// LRUCacheV2 implements the LRU cache class that implements the Cache and MCache interfaces.
//
// Key and value support generics.
// When the number of caches exceeds the limit, the contents of the cache are eliminated according to the usage of the Key, that is, the least used will be eliminated.
// The expiration time of the data does not affect the elimination strategy.

type LRUCacheV2[K comparable, V any] struct {
	lruMap  map[any]*item2[V]
	lruList *list.List
	mux     sync.Mutex
	maxUsed int // Maximum usage count, value 0 means no limitation
	caption int // Number of cached items
}

// Get reads the content from the cache.
// Return value explanation:
//
//	1st: cache value
//	2nd: cache existence, when true, the first parameter is valid
//	3rd: error information

func (lc *LRUCacheV2[K, V]) Get(_ context.Context, key K) (V, bool, error) {
	lc.mux.Lock()
	defer lc.mux.Unlock()

	val, ok := lc.getOne(key)
	return val, ok, nil
}

func (lc *LRUCacheV2[K, V]) getOne(key K) (V, bool) {
	tmp, ok := lc.lruMap[key]
	if !ok {
		var emp V
		return emp, false
	}
	// Check age
	if lc.maxUsed > 0 {
		if tmp.usedCount >= lc.maxUsed {
			delete(lc.lruMap, key)
			lc.lruList.Remove(tmp.el)
			var emp V
			return emp, false
		}
		tmp.usedCount++
	}
	// Check expiration
	if tmp.expireTime.Before(time.Now()) {
		delete(lc.lruMap, key)
		lc.lruList.Remove(tmp.el)
		var emp V
		return emp, false
	}
	// Hit
	lc.lruList.MoveToFront(tmp.el)
	return tmp.val, true
}

// MGet reads multiple contents from the cache.
//
// Return value explanation:
//
//	1st: cache values
//	2nd: cache existence, when true, the first return value is valid
//	3rd: error information

func (lc *LRUCacheV2[K, V]) MGet(_ context.Context, keys ...K) ([]V, []bool, error) {
	if len(keys) == 0 {
		return nil, nil, nil
	}
	values := make([]V, len(keys))
	oks := make([]bool, len(keys))
	for idx, key := range keys {
		val, ok := lc.getOne(key)
		values[idx] = val
		oks[idx] = ok
	}
	return values, oks, nil
}

// Set writes to the cache and sets the expiration time to ttl.

func (lc *LRUCacheV2[K, V]) Set(_ context.Context, key K, value V, ttl time.Duration) error {
	lc.mux.Lock()
	defer lc.mux.Unlock()

	lc.doSet(key, value, ttl)

	return nil
}

func (lc *LRUCacheV2[K, V]) doSet(key K, value V, ttl time.Duration) {
	if tmp, ok := lc.lruMap[key]; ok {
		tmp.val = value
		tmp.usedCount = 0
		tmp.expireTime = time.Now().Add(ttl)
		lc.lruList.MoveToFront(tmp.el)
		return
	}

	el := lc.lruList.PushFront(key)
	lc.lruMap[key] = &item2[V]{
		val:        value,
		usedCount:  0,
		expireTime: time.Now().Add(ttl),
		el:         el,
	}

	for lc.lruList.Len() > lc.caption {
		last := lc.lruList.Back()
		delete(lc.lruMap, last.Value)
		lc.lruList.Remove(last)
	}
}

// MSet writes multiple entries to the cache and sets the expiration time to ttl.

func (lc *LRUCacheV2[K, V]) MSet(_ context.Context, kvs map[K]V, ttl time.Duration) error {
	lc.mux.Lock()
	defer lc.mux.Unlock()
	for key, val := range kvs {
		lc.doSet(key, val, ttl)
	}
	return nil
}

// Delete deletes multiple cache keys.

func (lc *LRUCacheV2[K, V]) Delete(_ context.Context, keys ...K) error {
	lc.mux.Lock()
	defer lc.mux.Unlock()

	for _, key := range keys {
		if v, ok := lc.lruMap[key]; ok {
			delete(lc.lruMap, key)
			lc.lruList.Remove(v.el)
		}
	}
	return nil
}

type item2[V any] struct {
	expireTime time.Time // Expiration time
	val        V
	el         *list.Element
	usedCount  int // Usage count
}
