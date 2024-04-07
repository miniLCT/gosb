# fastrand

`fastrand` is the fastest pseudo-random number generator in Go. Support most common APIs of `math/rand`.

This generator base on the Go runtime per-M structure, and the init-seed provided by the Go runtime, which means you can't add your seed, but these methods scale very well on multiple cores.

## Compare to math/rand

- **2 ~ 200x faster**
- Scales well on multiple cores
- **Not** provide a stable value stream (can't inject init-seed)
- Fix bugs in math/rand `Float64` and `Float32`  (since no need to preserve the value stream)


## benchmark
Go version: go1.19.8 darwin/amd64

CPU: VirtualApple @ 2.50GHz

OS: Apple M1

MEMORY: 16GB
```
name                       old time/op  new time/op  delta
SingleCore/Uint32()-8      13.8ns ± 1%   2.6ns ± 2%  -80.93%  (p=0.008 n=5+5)
SingleCore/Uint64()-8      13.9ns ± 0%   2.3ns ± 1%  -83.63%  (p=0.008 n=5+5)
SingleCore/Int()-8         15.1ns ±15%   5.2ns ± 1%  -65.74%  (p=0.016 n=5+4)
SingleCore/Intn(32)-8      14.2ns ± 4%  12.2ns ± 0%  -14.56%  (p=0.008 n=5+5)
SingleCore/Read/1024-8      590ns ± 3%   149ns ± 0%  -74.69%  (p=0.008 n=5+5)
SingleCore/Read/10240-8    5.73µs ± 3%  1.40µs ± 0%  -75.65%  (p=0.008 n=5+5)
SingleCore/Perm/1024-8     15.7µs ±11%   6.4µs ± 4%  -59.12%  (p=0.008 n=5+5)
SingleCore/Shuffle/1024-8  14.4µs ± 1%   9.1µs ± 6%  -36.78%  (p=0.008 n=5+5)
MultipleCore/Uint32()-8     121ns ± 2%     1ns ± 6%  -99.55%  (p=0.008 n=5+5)
MultipleCore/Uint64()-8     119ns ± 4%     1ns ±18%  -99.46%  (p=0.008 n=5+5)
MultipleCore/Int()-8        117ns ± 2%     1ns ± 1%  -99.17%  (p=0.008 n=5+5)
MultipleCore/Intn(32)-8     120ns ± 5%     1ns ± 1%  -98.87%  (p=0.008 n=5+5)
MultipleCore/Read/1024-8    703ns ± 1%   123ns ±10%  -82.51%  (p=0.008 n=5+5)
MultipleCore/Read/10240-8  6.37µs ± 1%  0.63µs ±20%  -90.17%  (p=0.008 n=5+5)
MultipleCore/Perm/1024-8    133µs ± 4%     4µs ±19%  -97.17%  (p=0.008 n=5+5)
```