# gosb
go silver bullet: a collection of go utilities.

# Go Required Version
need go 1.27 or later. ⚠️

# Install

go get github.com/miniLCT/gosb

# Highlights

- **range-over-func iterators** (`iter.Seq` / `iter.Seq2`, go1.23): every container can be
  ranged over, and composed with `slices.Collect` / `maps.Collect`:

  ```go
  tree := rbt.New[int, string](func(a, b int) bool { return a < b })
  for k, v := range tree.All() { // ascending, Backward() for descending
      _ = k
      _ = v
  }

  set := gset.NewSetFrom(slices.Values([]int{3, 1, 2}))
  for v := range set.All() {
      _ = v
  }
  ```

  Available on `glist.List` (`All`/`Backward`/`Values`), `gring.Ring`, `gset.Set`,
  `gqueue.Queue` (`All`/`Drain`), `gheap.Heap`, `rbt.RbTree`, `gskiplist.SkipList`,
  `gconcurrent/sync.Map`, `bitset.BitSet` and `cachex.MGetResult`/`LRUCacheV2`.

- **standard library first** (go1.21+): `gslice`, `gmap`, `gset` and `gnumerics` are now
  implemented on top of `slices`, `maps`, `cmp`, the built-in `min`/`max` and `clear`,
  while `hack/fastrand` keeps the fastest runtime-backed generator.
- **generic methods** (go1.27): e.g. `(*glist.List[T]).Transform[U any](func(T) U) *List[U]`.
- **modern concurrency** (go1.25): `sync.WaitGroup.Go` in `gconcurrent/mr`, plus
  `atomic.Bool` and `atomic.Pointer[error]` instead of `int32` CAS and `atomic.Value`.

# Naming

picked from the following list of names:

`GoLib`

`GoUtil` 

`GoTools`

`GoKitBox`

`GoSnippets`

`GoCrafts`

`GoForge`

`GoSwissKnife`

`GoNexus`

`GoGizmo`

