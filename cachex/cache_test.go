package cachex

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestChain_keysNotFound(t *testing.T) {
	cc := &Chain[string, string]{}
	got1 := cc.keysNotFound([]string{"k1", "k2", "k3"}, []bool{true, false, false})
	want1 := []string{"k2", "k3"}
	require.Equal(t, want1, got1)

	got2 := cc.keysNotFound([]string{"k1", "k2", "k3"}, nil)
	want2 := []string{"k1", "k2", "k3"}
	require.Equal(t, want2, got2)

	got3 := cc.keysNotFound([]string{"k1", "k2", "k3"}, []bool{true, false})
	want3 := []string{"k2", "k3"}
	require.Equal(t, want3, got3)
}

func testGetOK(t *testing.T, key string, value string, caches ...Getter[string, string]) {
	for _, cache := range caches {
		got1, ok1, err1 := cache.Get(context.Background(), key)
		require.NoError(t, err1)
		require.Equal(t, value, got1)
		require.True(t, ok1)
	}
}

func testGetNot(t *testing.T, key string, caches ...Getter[string, string]) {
	for _, cache := range caches {
		got1, ok1, err1 := cache.Get(context.Background(), key)
		require.NoError(t, err1)
		require.Equal(t, "", got1)
		require.False(t, ok1)
	}
}

func testGetErr(t *testing.T, key string, caches ...Getter[string, string]) {
	for _, cache := range caches {
		got1, ok1, err1 := cache.Get(context.Background(), key)
		require.Error(t, err1)
		require.Equal(t, "", got1)
		require.False(t, ok1)
	}
}

func testMGetErr(t *testing.T, keys []string, caches ...MGetter[string, string]) {
	for _, cache := range caches {
		got1, ok1, err1 := cache.MGet(context.Background(), keys...)
		require.Error(t, err1)
		require.Nil(t, got1)
		require.Nil(t, ok1)
	}
}

func testMGetOk(t *testing.T, keys []string, wantValues []string, wantStatus []bool, caches ...MGetter[string, string]) {
	for _, cache := range caches {
		got1, ok1, err1 := cache.MGet(context.Background(), keys...)
		require.NoError(t, err1)
		require.Equal(t, wantValues, got1)
		require.Equal(t, wantStatus, ok1)
	}
}

func TestChain(t *testing.T) {
	c1 := NewLRUCacheV2[string, string](1000, 0)
	c2 := NewLRUCacheV2[string, string](1000, 0)
	cc := &Chain[string, string]{
		Caches: []*ChainItem[string, string]{
			{
				Cache: c1,
				TTL:   time.Second,
			},
			{
				Cache: c2,
				TTL:   time.Second,
			},
		},
	}

	testGetNot(t, "k1", c1, c2, cc, c1, c2)

	require.NoError(t, cc.Set(context.Background(), "k1", "v1", time.Minute))
	testGetOK(t, "k1", "v1", cc, c1, c2)

	require.NoError(t, c1.Delete(context.Background(), "k1"))
	testGetNot(t, "k1", c1)

	testGetOK(t, "k1", "v1", cc, c1)

	for i := 0; i < 3; i++ {
		vs1, st1, err1 := cc.MGet(context.Background(), "k1", "k2")
		require.NoError(t, err1)
		require.Equal(t, []string{"v1", ""}, vs1)
		require.Equal(t, []bool{true, false}, st1)
	}

	require.NoError(t, c2.Set(context.Background(), "k2", "v2", time.Minute))

	for i := 0; i < 3; i++ {
		testGetOK(t, "k2", "v2", c2)
		vs1, st1, err1 := cc.MGet(context.Background(), "k1", "k2")
		require.NoError(t, err1)
		require.Equal(t, []string{"v1", "v2"}, vs1)
		require.Equal(t, []bool{true, true}, st1)
		testGetOK(t, "k2", "v2", c1)
	}

	require.NoError(t, cc.Delete(context.Background(), "k1", "k2"))
	testGetNot(t, "k1", cc, c1, c2)
	testGetNot(t, "k2", cc, c1, c2)

	require.NoError(t, cc.MSet(context.Background(), nil, time.Second))
	require.NoError(t, cc.Delete(context.Background()))

	kv1 := map[string]string{
		"k10": "v10",
		"k11": "v11",
		"k12": "v12",
	}
	require.NoError(t, cc.MSet(context.Background(), kv1, time.Second))
	for k, v := range kv1 {
		testGetOK(t, k, v, c1, c2, cc)
	}

	vs2, st2, err2 := cc.MGet(context.Background())
	require.Nil(t, vs2)
	require.Nil(t, st2)
	require.NoError(t, err2)
}

func TestChain_1thErr(t *testing.T) {
	c1 := &testCache1[string, string]{}
	c2 := NewLRUCacheV2[string, string](1000, 0)
	cc := &Chain[string, string]{
		Caches: []*ChainItem[string, string]{
			{
				Cache: c1,
				TTL:   time.Second,
			},
			{
				Cache: c2,
				TTL:   time.Second,
			},
		},
	}

	testGetErr(t, "k1", cc, c1)

	// c1 Set 失败，但是 c2 应该 Set 成功
	require.Error(t, cc.Set(context.Background(), "k100", "v100", time.Minute))
	testGetOK(t, "k100", "v100", c2)

	// c1 删除失败，但是 c2 应该删除成功
	require.Error(t, cc.Delete(context.Background(), "k100", "k101"))
	testGetNot(t, "k100", c2)

	kv1 := map[string]string{
		"k101": "v101",
		"k102": "v102",
	}
	require.Error(t, cc.MSet(context.Background(), kv1, time.Minute))
	keys1 := []string{"k101", "k102"}
	values1 := []string{"v101", "v102"}
	status1 := []bool{true, true}
	testMGetErr(t, keys1, cc, c1)
	testMGetOk(t, keys1, values1, status1, c2)

	// ------------------
	cc.ContinueOnReadErr = true
	testGetNot(t, "k1", c2, cc)
	require.NoError(t, c2.Set(context.Background(), "k1", "v1", time.Minute))
	require.NoError(t, c2.Set(context.Background(), "k2", "v2", time.Minute))
	testGetOK(t, "k1", "v1", cc, c2)
	testGetOK(t, "k2", "v2", cc, c2)

	testMGetOk(t, keys1, values1, status1, cc, c2)
}

func TestNoCache(t *testing.T) {
	cc := &NoCache[string, string]{}
	require.NoError(t, cc.Delete(context.Background()))
	require.NoError(t, cc.Set(context.Background(), "k1", "v1", time.Second))
	require.NoError(t, cc.MSet(context.Background(), nil, time.Second))

	val1, ok1, err1 := cc.Get(context.Background(), "k1")
	require.NoError(t, err1)
	require.False(t, ok1)
	require.Equal(t, "", val1)

	vs2, ok2, err2 := cc.MGet(context.Background(), "k1", "k2")
	require.NoError(t, err2)
	require.Equal(t, []string{"", ""}, vs2)
	require.Equal(t, []bool{false, false}, ok2)

	vs3, ok3, err3 := cc.MGet(context.Background())
	require.NoError(t, err3)
	require.Nil(t, vs3)
	require.Nil(t, ok3)
}

var _ Cache[string, string] = (*testCache1[string, string])(nil)

type testCache1[K comparable, V any] struct {
	OnGet    func(ctx context.Context, key K) (V, bool, error)
	OnSet    func(ctx context.Context, key K, value V, ttl time.Duration) error
	OnDelete func(ctx context.Context, keys ...K) error
	OnMGet   func(ctx context.Context, keys ...K) ([]V, []bool, error)
	OnMSet   func(ctx context.Context, kvs map[K]V, ttl time.Duration) error
}

func (t *testCache1[K, V]) Get(ctx context.Context, key K) (V, bool, error) {
	if t.OnGet == nil {
		var emp V
		return emp, false, errors.New("no OnGet")
	}
	return t.OnGet(ctx, key)
}

func (t *testCache1[K, V]) Set(ctx context.Context, key K, value V, ttl time.Duration) error {
	if t.OnSet == nil {
		return errors.New("no OnSet")
	}
	return t.OnSet(ctx, key, value, ttl)
}

func (t *testCache1[K, V]) Delete(ctx context.Context, keys ...K) error {
	if t.OnDelete == nil {
		return errors.New("no OnDelete")
	}
	return t.OnDelete(ctx, keys...)
}

func (t *testCache1[K, V]) MGet(ctx context.Context, keys ...K) ([]V, []bool, error) {
	if t.OnMGet == nil {
		return nil, nil, errors.New("no OnMGet")
	}
	return t.OnMGet(ctx, keys...)
}

func (t *testCache1[K, V]) MSet(ctx context.Context, kvs map[K]V, ttl time.Duration) error {
	if t.OnMSet == nil {
		return errors.New("no OnMSet")
	}
	return t.OnMSet(ctx, kvs, ttl)
}

func TestMGet(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		c := &NoCache[string, string]{}
		ret, err := MGet[string, string](context.Background(), c, "k1")
		require.NoError(t, err)
		require.Equal(t, []string{"k1"}, ret.MissKeys())
		require.Empty(t, ret.HitKeys())
	})
}

func TestMGetResult(t *testing.T) {
	ret := MGetResult[string, string]{
		keys:   []string{"k1", "k2", "k3", "k4"},
		values: []string{"v1", "v2", "", ""},
		hits:   []bool{true, true, false, false},
	}
	hitKeys := []string{"k1", "k2"}
	require.Equal(t, hitKeys, ret.HitKeys())
	missKeys := []string{"k3", "k4"}
	require.Equal(t, missKeys, ret.MissKeys())

	hs, ms := ret.HitMissKeys()
	require.Equal(t, hitKeys, hs)
	require.Equal(t, missKeys, ms)

	var ks1 []string
	var vs1 []string
	ret.RangeHit(func(key string, Value string) bool {
		ks1 = append(ks1, key)
		vs1 = append(vs1, Value)
		return true
	})
	require.Equal(t, hitKeys, ks1)
	require.Equal(t, []string{"v1", "v2"}, vs1)

	var ks2 []string
	ret.RangeMiss(func(key string) bool {
		ks2 = append(ks2, key)
		return true
	})
	require.Equal(t, missKeys, ks2)

	var ks3 []string
	ret.Range(func(key string, hit bool, _ string) bool {
		ks3 = append(ks3, key)
		return false
	})
	require.Equal(t, []string{"k1"}, ks3)

	var ks4 []string
	ret.RangeHit(func(key string, _ string) bool {
		ks4 = append(ks4, key)
		return false
	})
	require.Equal(t, []string{"k1"}, ks4)

	var ks5 []string
	ret.RangeMiss(func(key string) bool {
		ks5 = append(ks5, key)
		return false
	})
	require.Equal(t, []string{"k3"}, ks5)

	g6 := ret.HitMapValues()
	want6 := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}
	require.Equal(t, want6, g6)
}
