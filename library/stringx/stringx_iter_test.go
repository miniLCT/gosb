package stringx

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLines(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// the yielded lines keep their terminating newline
	assert.Equal([]string{"a\n", "b\n"}, slices.Collect(Lines("a\nb\n")))
	assert.Equal([]string{"a\n", "b"}, slices.Collect(Lines("a\nb")))
}

func TestSplitSeq(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]string{"a", "b", "c"}, slices.Collect(SplitSeq("a,b,c", ",")))
	assert.Equal([]string{"abc"}, slices.Collect(SplitSeq("abc", ",")))
}

func TestFieldsSeq(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal([]string{"a", "b"}, slices.Collect(FieldsSeq(" a  b ")))
}

func TestSplitSeqBreak(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var got []string
	for v := range SplitSeq("a,b,c", ",") {
		got = append(got, v)
		if len(got) == 2 {
			break
		}
	}
	assert.Equal([]string{"a", "b"}, got)
}
