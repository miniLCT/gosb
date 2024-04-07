package fastrand

import (
	"math/rand"
	"os"
	"testing"
)

type testFuncs struct {
	fu32     func() uint32
	fu64     func() uint64
	fInt     func() int
	fIntn    func(_ int) int
	fRead    func(_ []byte) (int, error)
	fPerm    func(_ int) []int
	fShuffle func(_ int, swap func(i, j int))
}

var (
	targetFs   testFuncs
	mathRandFs = testFuncs{
		fu32:     rand.Uint32,
		fu64:     rand.Uint64,
		fInt:     rand.Int,
		fIntn:    rand.Intn,
		fRead:    rand.Read,
		fPerm:    rand.Perm,
		fShuffle: rand.Shuffle,
	}
	fastRandFs = testFuncs{
		fu32:     Uint32,
		fu64:     Uint64,
		fInt:     Int,
		fIntn:    Intn,
		fRead:    Read,
		fPerm:    Perm,
		fShuffle: Shuffle,
	}
)

func init() {
	if os.Getenv("BENCHMARK_TARGET") == "math_rand" {
		targetFs = mathRandFs
	} else {
		targetFs = fastRandFs
	}
}

func benchmarkSingleCore(b *testing.B, fs testFuncs) {
	b.Run("Uint32()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fs.fu32()
		}
	})
	b.Run("Uint64()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fs.fu64()
		}
	})
	b.Run("Int()", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fs.fInt()
		}
	})
	b.Run("Intn(32)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fs.fIntn(32)
		}
	})
	b.Run("Read/1024", func(b *testing.B) {
		p := make([]byte, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = fs.fRead(p)
		}
	})
	b.Run("Read/10240", func(b *testing.B) {
		p := make([]byte, 10240)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = fs.fRead(p)
		}
	})
	b.Run("Perm/1024", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fs.fPerm(1024)
		}
	})
	b.Run("Shuffle/1024", func(b *testing.B) {
		x := make([]int, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fs.fShuffle(1024, func(i, j int) {
				x[i], x[j] = x[j], x[i]
			})
		}
	})
}

func benchmarkMultipleCore(b *testing.B, fs testFuncs) {
	b.Run("Uint32()", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = fs.fu32()
			}
		})
	})
	b.Run("Uint64()", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = fs.fu64()
			}
		})
	})
	b.Run("Int()", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = fs.fInt()
			}
		})
	})
	b.Run("Intn(32)", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = fs.fIntn(32)
			}
		})
	})
	b.Run("Read/1024", func(b *testing.B) {
		p := make([]byte, 1024)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = fs.fRead(p)
			}
		})
	})
	b.Run("Read/10240", func(b *testing.B) {
		p := make([]byte, 10240)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = fs.fRead(p)
			}
		})
	})
	b.Run("Perm/1024", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fs.fPerm(1024)
			}
		})
	})
}

func BenchmarkSingleCore(b *testing.B) {
	benchmarkSingleCore(b, targetFs)
}

func BenchmarkMultipleCore(b *testing.B) {
	benchmarkMultipleCore(b, targetFs)
}
