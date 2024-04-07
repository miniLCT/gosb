gskiplist is the fastest implementation of skiplist in go which support generic.

#### Papers:

- [Skip Lists: A Probabilistic Alternative to Balanced Trees](https://www.cl.cam.ac.uk/teaching/0506/Algorithms/skiplists.pdf)

#### Functions:

Other than the classic `Find`, `Insert` and `Delete`, some more convenience functions are implemented that makes this skiplist implementation very easy and straight forward to use
in real applications. All complexity values are approximates, as skiplist can only approximate runtime complexity.

| Function        | Complexity           | Description  |
| ------------- |:-------------:|:-----|
| Find | O(log(n)) | Finds an element in the skiplist |
| FindGreaterOrEqual | O(log(n)) | Finds the first element that is greater or equal the given value in the skiplist |
| Insert | O(log(n)) | Inserts an element into the skiplist |
| Delete | O(log(n)) | Deletes an element from the skiplist |
| GetSmallestNode | O(1) | Returns the smallest element in the skiplist |
| GetLargestNode | O(1) | Returns the largest element in the skiplist |
| Prev | O(1) | Given a skiplist-node, it returns the previous element (Wraps around and allows to linearly iterate the skiplist) |
| Next | O(1) | Given a skiplist-node, it returns the next element (Wraps around and allows to linearly iterate the skiplist) |
| ChangeValue | O(1) | Given a skiplist-node, the actual value can be changed, as long as the key stays the same (Example: Change a structs data) |

#### Bench:

```text
$ go test -run=NOTEST -bench=. -count=1 -timeout=10m

goos: darwin
goarch: amd64
pkg: icode.baidu.com/baidu/passport/gogenerics/gds/gskiplist
cpu: VirtualApple @ 2.50GHz
BenchmarkLoadMostlyHits/*gskiplist.DeepCopyMap-8                330280492                3.658 ns/op
BenchmarkLoadMostlyHits/*gskiplist.RWMutexMap-8                 13389918                89.32 ns/op
BenchmarkLoadMostlyHits/*gskiplist.SyncMap[int,int]-8           159714217                7.327 ns/op
BenchmarkLoadMostlyHits/*gskiplist.SkipList[int,int]-8           8307650               145.4 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.DeepCopyMap-8              615504688                1.844 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.RWMutexMap-8               14186954                81.36 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.SyncMap[int,int]-8         364721040                3.372 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.SkipList[int,int]-8        874448346                1.482 ns/op
BenchmarkLoadOrStoreUnique/*gskiplist.RWMutexMap-8               3651057               339.9 ns/op
BenchmarkLoadOrStoreUnique/*gskiplist.SyncMap[int,int]-8         1000000              1265 ns/op
BenchmarkLoadOrStoreUnique/*gskiplist.SkipList[int,int]-8                 238684            284895 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.DeepCopyMap-8                   4190437               287.9 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.RWMutexMap-8                    8657164               138.2 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.SyncMap[int,int]-8              7297562               163.6 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.SkipList[int,int]-8            64378400                18.23 ns/op
BenchmarkAdversarialAlloc/*gskiplist.DeepCopyMap-8                       3414529               347.9 ns/op
BenchmarkAdversarialAlloc/*gskiplist.RWMutexMap-8                       23436660                73.72 ns/op
BenchmarkAdversarialAlloc/*gskiplist.SyncMap[int,int]-8                  6630920               218.9 ns/op
BenchmarkAdversarialAlloc/*gskiplist.SkipList[int,int]-8                 2665887              1011 ns/op
BenchmarkDeleteCollision/*gskiplist.DeepCopyMap-8                        8580079               143.4 ns/op
BenchmarkDeleteCollision/*gskiplist.RWMutexMap-8                        10476196               115.5 ns/op
BenchmarkDeleteCollision/*gskiplist.SyncMap[int,int]-8                  438461808                2.922 ns/op
BenchmarkDeleteCollision/*gskiplist.SkipList[int,int]-8                 1000000000               0.9350 ns/op
```
