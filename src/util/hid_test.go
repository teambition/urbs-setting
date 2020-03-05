package util

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHID(t *testing.T) {
	t.Run("HID should work", func(t *testing.T) {
		assert := assert.New(t)
		maxInt64 := int64(math.MaxInt64)

		hid1 := NewHID([]byte("abc"))
		hid2 := NewHID([]byte("123"))

		assert.Equal("", hid1.ToHex(-999))
		assert.Equal("", hid1.ToHex(-1))
		assert.Equal("", hid1.ToHex(0))
		// assert.Equal("", hid1.ToHex(maxInt64+1))

		s := hid1.ToHex(1)
		assert.Equal(int64(1), hid1.ToInt64(s))
		assert.Equal(int64(0), hid2.ToInt64(s))
		assert.Equal(24, len(s))

		s = hid1.ToHex(999)
		assert.Equal(int64(999), hid1.ToInt64(s))
		assert.Equal(int64(0), hid2.ToInt64(s))
		assert.Equal(24, len(s))

		s = hid1.ToHex(maxInt64)
		assert.Equal(maxInt64, hid1.ToInt64(s))
		assert.Equal(int64(0), hid2.ToInt64(s))
		assert.Equal(24, len(s))

		assert.NotEqual(hid1.ToHex(1), hid2.ToHex(1))
		assert.NotEqual(hid1.ToHex(maxInt64), hid2.ToHex(maxInt64))
	})
}
