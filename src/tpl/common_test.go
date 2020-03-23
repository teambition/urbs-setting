package tpl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToSlice(t *testing.T) {
	t.Run(`StringToSlice should work`, func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal([]string{}, StringToSlice(""))

		assert.Equal([]string{"a"}, StringToSlice("a"))
		assert.Equal([]string{"a", "b"}, StringToSlice("a,b"))
	})
}

func TestStringSliceHas(t *testing.T) {
	t.Run(`StringSliceHas should work`, func(t *testing.T) {
		assert := assert.New(t)
		assert.True(StringSliceHas([]string{"a", "b"}, "b"))
		assert.False(StringSliceHas([]string{"a", "bb"}, "b"))
		assert.False(StringSliceHas([]string{}, "b"))
		assert.False(StringSliceHas([]string{}, ""))
		assert.False(StringSliceHas([]string{"a"}, ""))
	})
}

func TestSortStringsAndCheck(t *testing.T) {
	t.Run(`SortStringsAndCheck should work`, func(t *testing.T) {
		assert := assert.New(t)

		s := []string{}
		assert.True(SortStringsAndCheck(s))

		s = []string{"a"}
		assert.True(SortStringsAndCheck(s))

		s = []string{""}
		assert.False(SortStringsAndCheck(s))

		s = []string{"b", "c", "a"}
		assert.True(SortStringsAndCheck(s))
		assert.Equal([]string{"a", "b", "c"}, s)

		s = []string{"b", "c", "a", "c"}
		assert.False(SortStringsAndCheck(s))

		s = []string{"b", "c", "a", ""}
		assert.False(SortStringsAndCheck(s))
	})
}
