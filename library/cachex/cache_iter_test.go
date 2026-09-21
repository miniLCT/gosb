package cachex

import (
	"context"
	"maps"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLRUCacheV2All(t *testing.T) {
	ctx := context.Background()
	cache := NewLRUCacheV2[string, int](10, 0)

	require.NoError(t, cache.Set(ctx, "a", 1, time.Minute))
	require.NoError(t, cache.Set(ctx, "b", 2, time.Minute))

	require.Equal(t, map[string]int{"a": 1, "b": 2}, maps.Collect(cache.All()))

	// the key keeps its generic type, no more map[any]
	m := maps.Collect(cache.All())
	for k := range m {
		require.IsType(t, "string", k)
	}
}

func TestLRUCacheV2AllIsEmptyAfterDelete(t *testing.T) {
	ctx := context.Background()
	cache := NewLRUCacheV2[string, int](10, 0)

	require.NoError(t, cache.Set(ctx, "a", 1, time.Minute))
	require.NoError(t, cache.Delete(ctx, "a"))

	var n int
	for range cache.All() {
		n++
	}
	require.Equal(t, 0, n)
}

func TestMGetResultIterators(t *testing.T) {
	res := MGetResult[string, int]{
		keys:   []string{"a", "b", "c"},
		values: []int{1, 0, 3},
		hits:   []bool{true, false, true},
	}

	require.Equal(t, map[string]int{"a": 1, "b": 0, "c": 3}, maps.Collect(res.All()))

	var hitKeys []string
	for k := range res.Hits() {
		hitKeys = append(hitKeys, k)
	}
	slices.Sort(hitKeys)
	require.Equal(t, []string{"a", "c"}, hitKeys)

	var missKeys []string
	for k := range res.Misses() {
		missKeys = append(missKeys, k)
	}
	require.Equal(t, []string{"b"}, missKeys)
}
