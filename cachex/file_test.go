package cachex

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	fc := &FileStore[string, string]{
		Dir:     t.TempDir(),
		GCCycle: time.Second,
	}
	defer func() {
		_ = fc.Purge()
	}()

	testGetNot(t, "key1", fc)
	require.NoError(t, fc.Set(context.Background(), "k1", "v1", time.Minute))
	checkOk := func(t *testing.T, key string, value string) {
		testGetOK(t, key, value, fc)
	}
	checkOk(t, "k1", "v1")
	require.NoError(t, fc.Delete(context.Background(), "k1", "k2"))
	testGetNot(t, "key1", fc)

	kv1 := map[string]string{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}
	require.NoError(t, fc.MSet(context.Background(), kv1, time.Second))
	checkOk(t, "k1", "v1")
	checkOk(t, "k2", "v2")
	checkOk(t, "k3", "v3")

	val2, ok2, err2 := fc.MGet(context.Background(), "k1", "k2", "k4")
	require.NoError(t, err2)
	require.Equal(t, []string{"v1", "v2", ""}, val2)
	require.Equal(t, []bool{true, true, false}, ok2)

	fp := fc.getFilePath("k1")
	_, err1 := os.Stat(fp)
	require.NoError(t, err1)

	time.Sleep(time.Second)
	require.NoError(t, fc.Delete(context.Background()))

	time.Sleep(time.Second) // Wait for background cleaning to complete
	_, err3 := os.Stat(fp)
	require.Error(t, err3)

	val4, _, err4 := fc.MGet(context.Background())
	require.NoError(t, err4)
	require.Nil(t, val4)

	gor1 := runtime.NumGoroutine()
	for i := 0; i < 1000; i++ {
		fc.gc()
	}
	require.NoError(t, fc.Set(context.Background(), "k1", "v1", time.Minute))
	require.NoError(t, fc.Purge())
	gor2 := runtime.NumGoroutine()
	require.LessOrEqual(t, gor2, gor1+2)
}

func TestFileStore1(t *testing.T) {
	fc := &FileStore[string, string]{}
	require.Error(t, fc.Delete(context.Background(), "k1"))
	require.Error(t, fc.Set(context.Background(), "k1", "v1", time.Second))
	require.Error(t, fc.MSet(context.Background(), map[string]string{"k1": "v1"}, time.Second))
	require.Error(t, fc.Purge())
}

func TestFileStore2(t *testing.T) {
	fc := &FileStore[string, string]{
		Dir:     t.TempDir(),
		GCCycle: time.Second,
	}
	require.NoError(t, fc.Set(context.Background(), "k100", "v100", time.Minute))
	testGetOK(t, "k100", "v100", fc)

	fp1 := fc.getFilePath("k100")

	_, err1 := os.Stat(fp1)
	require.NoError(t, err1)

	require.NoError(t, os.WriteFile(fp1, []byte("hello"), 0644))
	testGetErr(t, "k100", fc)
	// After the file content is incorrect, it will be automatically deleted
	_, err2 := os.Stat(fp1)
	require.Error(t, err2)
}
