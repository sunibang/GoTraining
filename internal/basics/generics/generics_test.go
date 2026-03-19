package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericStack(t *testing.T) {
	// Stack of ints
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)

	val, ok := intStack.Pop()
	assert.True(t, ok)
	assert.Equal(t, 20, val)

	// Stack of strings
	stringStack := &Stack[string]{}
	stringStack.Push("hello")
	valStr, ok := stringStack.Pop()
	assert.True(t, ok)
	assert.Equal(t, "hello", valStr)
}

func TestSliceContains(t *testing.T) {
	assert.True(t, SliceContains([]int{1, 2, 3}, 2))
	assert.True(t, SliceContains([]string{"a", "b"}, "a"))
	assert.False(t, SliceContains([]float64{1.1, 2.2}, 3.3))
}

func TestSumWithUnderlyingTypes(t *testing.T) {
	type MyInt int
	vals := []MyInt{1, 2, 3}

	// This works because of the '~' in our Number constraint!
	total := Sum(vals)
	assert.Equal(t, MyInt(6), total)
}

func TestMap(t *testing.T) {
	nums := []int{1, 2, 3}
	strs := Map(nums, func(n int) string {
		return "num" + strconv.Itoa(n)
	})

	assert.Equal(t, []string{"num1", "num2", "num3"}, strs)
}
