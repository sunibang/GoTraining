package casting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumericConversion(t *testing.T) {
	// Go is strict: NO implicit conversion.
	var i = 42
	var f = float64(i) // Must convert manually
	var u = uint(f)    // Must convert manually

	assert.Equal(t, 42, i)
	assert.Equal(t, 42.0, f)
	assert.Equal(t, uint(42), u)
}

func TestLossyConversion(t *testing.T) {
	// Be careful with precision loss!
	var f = 42.99
	var i = int(f) // Truncates towards zero

	assert.Equal(t, 42, i) // .99 is lost
}

func TestOverflowConversion(t *testing.T) {
	// Be careful with overflow!
	var bigInt int64 = 257
	var smallInt = int8(bigInt)

	// 257 in binary is 1 0000 0001
	// int8 only takes the last 8 bits: 0000 0001 = 1
	assert.Equal(t, int8(1), smallInt)
}

func TestStringConversion(t *testing.T) {
	// Converting between string and byte/rune slices
	s := "hello"
	b := []byte(s)
	r := []rune(s)

	assert.Equal(t, []byte{'h', 'e', 'l', 'l', 'o'}, b)
	assert.Equal(t, s, string(b))
	assert.Equal(t, []rune{'h', 'e', 'l', 'l', 'o'}, r)
}
