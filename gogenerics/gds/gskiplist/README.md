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
pkg: github.com/miniLCT/gosb/gds/gskiplist
cpu: VirtualApple @ 2.50GHz
BenchmarkLoadMostlyHits/*gskiplist.DeepCopyMap-8                339300584                3.648 ns/op
BenchmarkLoadMostlyHits/*gskiplist.RWMutexMap-8                 15707191                82.62 ns/op
BenchmarkLoadMostlyHits/*gskiplist.SyncMap[int,int]-8           148558563                8.203 ns/op
BenchmarkLoadMostlyHits/*gskiplist.SkipList[int,int]-8           8557464               144.6 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.DeepCopyMap-8              691986650                1.965 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.RWMutexMap-8               16790496                89.97 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.SyncMap[int,int]-8         252446268                4.683 ns/op
BenchmarkLoadMostlyMisses/*gskiplist.SkipList[int,int]-8        814085949                1.394 ns/op
BenchmarkLoadOrStoreUnique/*gskiplist.RWMutexMap-8               3293228               347.4 ns/op
BenchmarkLoadOrStoreUnique/*gskiplist.SyncMap[int,int]-8         1890684               782.6 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.DeepCopyMap-8                   3802148               326.7 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.RWMutexMap-8                    8502034               140.1 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.SyncMap[int,int]-8              7065519               169.8 ns/op
BenchmarkLoadOrStoreCollision/*gskiplist.SkipList[int,int]-8            64089868                18.39 ns/op
BenchmarkAdversarialAlloc/*gskiplist.DeepCopyMap-8                       2865942               393.4 ns/op
BenchmarkAdversarialAlloc/*gskiplist.RWMutexMap-8                       22905480                74.07 ns/op
BenchmarkAdversarialAlloc/*gskiplist.SyncMap[int,int]-8                  6309598               212.4 ns/op
BenchmarkAdversarialAlloc/*gskiplist.SkipList[int,int]-8                 2038196               932.5 ns/op
BenchmarkDeleteCollision/*gskiplist.DeepCopyMap-8                        6336438               180.4 ns/op
BenchmarkDeleteCollision/*gskiplist.RWMutexMap-8                        10384668               117.9 ns/op
BenchmarkDeleteCollision/*gskiplist.SyncMap[int,int]-8                  267763987                4.468 ns/op
BenchmarkDeleteCollision/*gskiplist.SkipList[int,int]-8                 779491996                1.544 ns/op
```