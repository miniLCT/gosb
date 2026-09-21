# gogenerics

This is a Go package that provides generic types and functions.

Requires go 1.27 or later.

## Packages

| Package | Description |
|---|---|
| `gconstraints` | constraints (`Ordered` is an alias of `cmp.Ordered`) and functional types |
| `gnumerics` | `Max`/`Min` family built on the built-in `min`/`max` and `slices` |
| `gcontainers/gslice` | slice helpers, delegates to `slices` + `iter` |
| `gcontainers/gmap` | map helpers, delegates to `maps` + `iter` |
| `gcontainers/gset` | generic hash set with `All()` iterator and `clear` |
| `gcontainers/glist` | generic doubly linked list with `All`/`Backward`/`Values` |
| `gcontainers/gring` | generic ring with `All`/`Backward` |
| `gcontainers/gqueue` | FIFO queue with `All` (non-destructive) and `Drain` |
| `gcontainers/gheap` | generic heap with `All`, `Peek` and `TryPop` |
| `gcontainers/gpool` | generic `sync.Pool` |
| `gconcurrent/sync` | generic `sync.Map` with `All` and `Clear` |
| `gconcurrent/mr` | mapreduce based on `sync.WaitGroup.Go` |
| `gds/grbt` | red black tree with `All`/`Backward`/`Min`/`Max` |
| `gds/gskiplist` | skiplist with `All`/`Backward` |

## Iteration

Every container exposes a `range`-able iterator, so it works with the `slices`
and `maps` helpers:

```go
l := glist.NewFrom(slices.Values([]int{1, 2, 3}))

values := slices.Collect(l.Values())          // []int{1, 2, 3}
doubled := l.Transform(func(v int) int { return v * 2 })
for v := range doubled.Values() {
    _ = v
}
```
