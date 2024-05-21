package cachex

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var _ Cache[string, any] = (*FileStore[string, any])(nil)

const cacheFileExt = ".cache"

// FileStore Local file cache
//
// Note:
//
//	Please make sure the usage scenario is appropriate to avoid occupying too much disk space or I/O exceptions.
//	It is recommended to have fewer than 10,000 keys.
//	A typical example is caching information for C-end users (too many keys, not enumerable), which is very inappropriate.

type FileStore[K comparable, V any] struct {
	// Dir Data root directory path, required
	Dir string

	// GCCycle Interval cycle for background cleaning of expired data, optional, default is 1 minute
	//
	// Every GCCycle duration, when calling Set, Get, and other interfaces, a background goroutine is started to scan files and delete expired cache files.
	GCCycle time.Duration

	lastGC   atomic.Int64
	gcStatus atomic.Bool
}

func (f *FileStore[K, V]) getFilePath(key K) string {
	str := fmt.Sprint(key)
	h := md5.New()
	h.Write([]byte(str))
	s := hex.EncodeToString(h.Sum(nil))
	fp := filepath.Join(f.Dir, s[:3], s[3:6], s[6:9], s[9:12], s[12:15], s[16:])
	return fp + cacheFileExt
}

// Get Read the content from the cache
// Return values:
//
//	1st: cache value
//	2nd: whether the cache exists, when true, the first parameter is valid
//	3rd: error information
func (f *FileStore[K, V]) Get(_ context.Context, key K) (V, bool, error) {
	defer f.gc()
	fp := f.getFilePath(key)
	return f.readFile(fp)
}

func (f *FileStore[K, V]) readFile(fp string) (V, bool, error) {
	content, err := os.ReadFile(fp)
	if err != nil {
		var emp V
		if os.IsNotExist(err) {
			return emp, false, nil
		}
		return emp, false, err
	}
	item := &fileCacheItem[K, V]{}
	err = json.Unmarshal(content, &item)
	if err != nil {
		_ = os.Remove(fp)
		var emp V
		return emp, false, err
	}
	if item.Alive() {
		return item.Data, true, nil
	}
	_ = os.Remove(fp)
	var emp V
	return emp, false, err
}

var errDirEmpty = errors.New("cache Dir is empty")

// Set Write to cache, and set the expiration time to ttl
func (f *FileStore[K, V]) Set(ctx context.Context, key K, value V, ttl time.Duration) error {
	if f.Dir == "" {
		return errDirEmpty
	}
	defer f.gc()

	fp := f.getFilePath(key)
	item := fileCacheItem[K, V]{
		Data: value,
		Key:  key,
		TTL:  int64(ttl),
		Exp:  time.Now().Add(ttl).UnixNano(),
	}
	content, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(fp)
	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0777)
	}
	return os.WriteFile(fp, content, 0666)
}

// MGet Batch read the content from the cache
//
// Return values:
//
//	1st: cache values
//	2nd: whether the caches exist, when true, the first return value is valid
//	3rd: error information
func (f *FileStore[K, V]) MGet(ctx context.Context, keys ...K) ([]V, []bool, error) {
	defer f.gc()
	if len(keys) == 0 {
		return nil, nil, nil
	}
	values := make([]V, len(keys))
	status := make([]bool, len(keys))
	for idx, key := range keys {
		val, st, err := f.Get(ctx, key)
		if err != nil {
			return nil, nil, err
		}
		values[idx] = val
		status[idx] = st
	}
	return values, status, nil
}

// MSet Batch write to cache, and set the expiration time to ttl
func (f *FileStore[K, V]) MSet(ctx context.Context, kvs map[K]V, ttl time.Duration) error {
	if f.Dir == "" {
		return errDirEmpty
	}

	defer f.gc()

	if len(kvs) == 0 {
		return nil
	}
	var errs []error
	for k, v := range kvs {
		if err := f.Set(ctx, k, v, ttl); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Delete Batch delete cache keys
//
// When multiple keys are provided, the corresponding cache files are deleted serially, and all errors are returned in the end.
func (f *FileStore[K, V]) Delete(ctx context.Context, keys ...K) error {
	if f.Dir == "" {
		return errDirEmpty
	}
	defer f.gc()
	if len(keys) == 0 {
		return nil
	}
	var errs []error
	for _, k := range keys {
		fp := f.getFilePath(k)
		err := os.Remove(fp)
		if err == nil || os.IsNotExist(err) {
			continue
		}
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

// gc Background cleaning of expired cache files, only one background goroutine performs cleaning within each cycle
func (f *FileStore[K, V]) gc() {
	if f.gcStatus.Load() {
		return
	}

	if f.lastGC.Load() == 0 {
		f.lastGC.Store(time.Now().Unix())
		return
	}

	cycle := f.GCCycle
	if cycle <= 0 {
		cycle = time.Minute
	}

	if time.Now().Unix()-f.lastGC.Load() < int64(cycle/time.Second) {
		return
	}
	if !f.gcStatus.CompareAndSwap(false, true) {
		return
	}
	go func() {
		defer func() {
			f.lastGC.Store(time.Now().Unix())
			f.gcStatus.Store(false)
			_ = recover()
		}()
		f.scanExpire()
	}()
}

func (f *FileStore[K, V]) scanExpire() {
	_ = filepath.WalkDir(f.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		_, _, _ = f.readFile(path)
		return nil
	})
}

// Purge Delete all cache content
func (f *FileStore[K, V]) Purge() error {
	if f.Dir == "" {
		return errDirEmpty
	}
	var lastErr error
	var errTotal int
	_ = filepath.WalkDir(f.Dir, func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, cacheFileExt) {
			e := os.Remove(path)
			if e != nil && !os.IsNotExist(err) {
				errTotal++
				lastErr = e
			}
		}
		return nil
	})
	if lastErr == nil {
		return nil
	}
	return fmt.Errorf("%d errors, last error is %w", errTotal, lastErr)
}

type fileCacheItem[K comparable, V any] struct {
	Data  V
	Key   K
	TTL   int64
	Ctime int64
	Exp   int64
}

func (fi *fileCacheItem[K, V]) Alive() bool {
	return fi.Exp > time.Now().UnixNano()
}
