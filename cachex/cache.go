package cachex

import (
	"context"
	"errors"
	"time"
)

// Cache is the interface definition of a cache object that supports generic Key, Value, single,
// and batch functionalities.

type Cache[K comparable, V any] interface {
	SCache[K, V]
	MCache[K, V]
}

// SCache is the interface definition of a cache object (single) that supports generic Key, Value.
type SCache[K comparable, V any] interface {
	Getter[K, V]
	Setter[K, V]
	Deleter[K, V]
}

type (
	// Getter is the cache read interface definition.
	Getter[K comparable, V any] interface {
		// Get reads the content from the cache.
		// Return values:
		//   1st: cache value
		//   2nd: whether cache exists, when true, the first parameter is valid
		//   3rd: error message
		Get(ctx context.Context, key K) (V, bool, error)
	}

	// Setter is the cache write interface definition.
	Setter[K comparable, V any] interface {
		// Set writes to the cache and sets the expiration time to ttl.
		Set(ctx context.Context, key K, value V, ttl time.Duration) error
	}
)

// MCache is the interface definition of a cache object (batch) that supports generic Key, Value.
type MCache[K comparable, V any] interface {
	MGetter[K, V]
	MSetter[K, V]
	Deleter[K, V]
}

type (
	// MGetter is the interface definition for batch querying cache data.
	MGetter[K comparable, V any] interface {
		// MGet reads the content from the cache.
		// Return values:
		//   1st: cache value
		//   2nd: whether cache exists, when true, the first return value is valid
		//   3rd: error message
		MGet(ctx context.Context, keys ...K) ([]V, []bool, error)
	}

	// MSetter is the interface definition for batch cache write.
	MSetter[K comparable, V any] interface {
		// MSet writes to the cache in bulk and sets the expiration time to ttl.
		MSet(ctx context.Context, kvs map[K]V, ttl time.Duration) error
	}
)

// Deleter is the cache deletion interface definition.
type Deleter[K comparable, V any] interface {
	// Delete deletes cache keys in batches, and if the key does not exist, it will also be considered deleted successfully.
	//
	// When deleting multiple keys at the same time, if some are successfully deleted and some fail, an error will also be returned (error!=nil).
	Delete(ctx context.Context, keys ...K) error
}

var _ Cache[string, any] = (*Chain[string, any])(nil)

type ChainItem[K comparable, V any] struct {
	// Cache is the cache object, required.
	Cache Cache[K, V]

	// TTL sets the expiration time for the cache, required.
	TTL time.Duration
}

// Chain is a multi-level cache that queries caches in sequence.
//
// If the upper-level cache does not return results, continue querying from the next cache and store the results in the upper-level cache.
// A typical application is a 2-level cache:
//  1. LRU Cache: with a short TTL
//  2. Redis Cache: with a longer TTL
//
// When LRU Cache has no results, continue querying from Redis Cache. If the result exists, it will be stored in LRU Cache.
// If there are still no results, it will return without caching automatically.
type Chain[K comparable, V any] struct {
	// Caches, required, list of multi-level caches.
	//
	// Query sequentially, until there is a result, and set the result to the Cache object that had no result before.
	Caches []*ChainItem[K, V]

	// ContinueOnReadErr continues to query the next cache if the current cache query returns an error, default is false.
	// This parameter currently only takes effect in GET and MGET methods.
	ContinueOnReadErr bool
}

// Get reads the content from the cache.
// Return values:
//
//	1st: cache value
//	2nd: whether cache exists, when true, the first parameter is valid
//	3rd: error message
func (c *Chain[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	var value V
	var has bool
	var err error
	for idx, item := range c.Caches {
		value, has, err = item.Cache.Get(ctx, key)
		if has {
			// Set the content queried from the next-level cache to the upper-level cache
			c.setForGet(ctx, c.Caches[:idx], key, value)
			return value, true, nil
		}
		if err != nil && !c.ContinueOnReadErr {
			break
		}
	}
	return value, has, err
}

func (c *Chain[K, V]) setForGet(ctx context.Context, caches []*ChainItem[K, V], key K, value V) {
	if len(caches) == 0 {
		return
	}
	for _, item := range caches {
		_ = item.Cache.Set(ctx, key, value, item.TTL)
	}
}

// Set sets all caches one by one and returns an error list.
//
// The TTL parameter passed to this method is invalid.
func (c *Chain[K, V]) Set(ctx context.Context, key K, value V, _ time.Duration) error {
	var errs []error
	for _, item := range c.Caches {
		err := item.Cache.Set(ctx, key, value, item.TTL)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (c *Chain[K, V]) MGet(ctx context.Context, keys ...K) ([]V, []bool, error) {
	if len(keys) == 0 {
		return nil, nil, nil
	}
	length := len(keys)
	indexMap := make(map[K]int, length)
	for i, k := range keys {
		indexMap[k] = i
	}
	values := make([]V, length)
	status := make([]bool, length)
	var err error

	for idx, item := range c.Caches {
		if len(keys) == 0 {
			break
		}
		vals, oks, err1 := item.Cache.MGet(ctx, keys...)
		err = err1
		if err1 != nil && !c.ContinueOnReadErr {
			return nil, nil, err1
		}
		for j := 0; j < len(oks); j++ {
			if oks[j] {
				keyIndex := indexMap[keys[j]]
				values[keyIndex] = vals[j]
				status[keyIndex] = true
			}
		}
		// Set the content queried from the next-level cache to the upper-level cache
		c.msetForMGet(ctx, c.Caches[:idx], keys, vals, oks)
		keys = c.keysNotFound(keys, oks)
	}
	return values, status, err
}

func (c *Chain[K, V]) msetForMGet(ctx context.Context, caches []*ChainItem[K, V], keys []K, values []V, status []bool) {
	if len(caches) == 0 || len(values) == 0 {
		return
	}
	kvs := make(map[K]V, len(keys))
	for idx, ok := range status {
		if ok {
			kvs[keys[idx]] = values[idx]
		}
	}
	if len(kvs) == 0 {
		return
	}
	for _, item := range caches {
		_ = item.Cache.MSet(ctx, kvs, item.TTL)
	}
}

// keysNotFound filters out the list of keys without results
func (c *Chain[K, V]) keysNotFound(keys []K, status []bool) []K {
	result := make([]K, 0, len(keys))
	for idx, key := range keys {
		if len(status) <= idx || !status[idx] {
			result = append(result, key)
		}
	}
	return result
}

// MSet sets all caches in bulk and returns an error list.
func (c *Chain[K, V]) MSet(ctx context.Context, kvs map[K]V, _ time.Duration) error {
	if len(kvs) == 0 {
		return nil
	}
	var errs []error
	for _, item := range c.Caches {
		err := item.Cache.MSet(ctx, kvs, item.TTL)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// Delete deletes all caches one by one and returns an error list.
func (c *Chain[K, V]) Delete(ctx context.Context, keys ...K) error {
	if len(keys) == 0 {
		return nil
	}
	var errs []error
	for _, item := range c.Caches {
		err := item.Cache.Delete(ctx, keys...)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

var _ Cache[string, any] = (*NoCache[string, any])(nil)

// NoCache is an empty cache, querying always returns no result, writing always succeeds.

type NoCache[K comparable, V any] struct{}

func (n *NoCache[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	var emp V
	return emp, false, nil
}

func (n *NoCache[K, V]) Set(ctx context.Context, key K, value V, ttl time.Duration) error {
	return nil
}

func (n *NoCache[K, V]) Delete(ctx context.Context, keys ...K) error {
	return nil
}

func (n *NoCache[K, V]) MGet(ctx context.Context, keys ...K) ([]V, []bool, error) {
	if len(keys) == 0 {
		return nil, nil, nil
	}
	values := make([]V, len(keys))
	status := make([]bool, len(keys))
	return values, status, nil
}

func (n *NoCache[K, V]) MSet(ctx context.Context, kvs map[K]V, ttl time.Duration) error {
	return nil
}

// MGet is an auxiliary method for MGetter, providing a friendlier result.
func MGet[K comparable, V any](ctx context.Context, cache MGetter[K, V], keys ...K) (MGetResult[K, V], error) {
	values, hits, err := cache.MGet(ctx, keys...)
	return MGetResult[K, V]{
		keys:   keys,
		values: values,
		hits:   hits,
	}, err
}

// MGetResult is the return value of the MGet method, providing a series of auxiliary functions for the result.

type MGetResult[K comparable, V any] struct {
	keys   []K
	values []V
	hits   []bool
}

// Range traverses all results.
func (mr MGetResult[K, V]) Range(fn func(key K, hit bool, value V) bool) {
	for idx := 0; idx < len(mr.keys); idx++ {
		if !fn(mr.keys[idx], mr.hits[idx], mr.values[idx]) {
			return
		}
	}
}

// RangeHit traverses all results hit by the cache.
func (mr MGetResult[K, V]) RangeHit(fn func(key K, Value V) bool) {
	for idx := 0; idx < len(mr.keys); idx++ {
		if !mr.hits[idx] {
			continue
		}
		if !fn(mr.keys[idx], mr.values[idx]) {
			return
		}
	}
}

// RangeMiss traverses all results missed by the cache.
func (mr MGetResult[K, V]) RangeMiss(fn func(key K) bool) {
	for idx := 0; idx < len(mr.keys); idx++ {
		if mr.hits[idx] {
			continue
		}
		if !fn(mr.keys[idx]) {
			return
		}
	}
}

// HitKeys returns a list of keys hit by the cache.
func (mr MGetResult[K, V]) HitKeys() []K {
	result := make([]K, 0, len(mr.keys))
	mr.RangeHit(func(key K, _ V) bool {
		result = append(result, key)
		return true
	})
	return result
}

// HitMapValues returns the map result data hit by the cache.
func (mr MGetResult[K, V]) HitMapValues() map[K]V {
	result := make(map[K]V, len(mr.keys))
	mr.RangeHit(func(key K, val V) bool {
		result[key] = val
		return true
	})
	return result
}

// MissKeys returns a list of keys missed by the cache.
func (mr MGetResult[K, V]) MissKeys() []K {
	result := make([]K, 0, len(mr.keys))
	mr.RangeMiss(func(key K) bool {
		result = append(result, key)
		return true
	})
	return result
}

// HitMissKeys returns lists of keys hit and missed by the cache separately.
func (mr MGetResult[K, V]) HitMissKeys() (hit []K, miss []K) {
	hit = make([]K, 0, len(mr.keys))
	miss = make([]K, 0, len(mr.keys))
	for idx := 0; idx < len(mr.keys); idx++ {
		if mr.hits[idx] {
			hit = append(hit, mr.keys[idx])
		} else {
			miss = append(miss, mr.keys[idx])
		}
	}
	return hit, miss
}
