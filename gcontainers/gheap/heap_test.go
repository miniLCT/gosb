package gheap

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/miniLCT/gosb/gcontainers/gslice"
)

var charSet = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "0",
}

func genString() string {
	l := 4
	var str string
	for i := 0; i < l; i++ {
		str += charSet[rand.Intn(len(charSet))]
	}
	return str
}

type Node struct {
	Name  string
	Count int
}

func TestHeap(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	rand.Seed(time.Now().Unix())
	nds := []*Node{}
	for i := 0; i < 10; i++ {
		nds = append(nds, &Node{
			Name:  genString(),
			Count: rand.Intn(4),
		})
	}
	less := func(a, b *Node) bool {
		if a.Count != b.Count {
			return a.Count > b.Count
		}
		return a.Name < b.Name
	}
	hp := NewWithData(nds, less)

	tmp := []*Node{}
	for i := 0; i <= 100; i++ {
		Push(hp, &Node{
			Name:  genString(),
			Count: rand.Intn(4),
		})
		for len(hp.data) > 0 {
			v := Pop(hp)
			tmp = append(tmp, v)
		}
		assert.True(gslice.IsSortedFunc(tmp, less))
		for _, v := range tmp {
			Push(hp, v)
		}
		tmp = make([]*Node, 0)
	}

	for i := 0; i <= 100; i++ {
		tmpSlice := gslice.Copy(hp.data)
		gslice.Shuffle(tmpSlice)
		h := NewWithData(tmpSlice[:rand.Int()%len(tmpSlice)], less)
		tpchecks := make([]*Node, 0)
		for len(h.data) > 0 {
			v := Pop(h)
			tpchecks = append(tpchecks, v)
		}
		assert.True(gslice.IsSortedFunc(tpchecks, less))
	}
}
