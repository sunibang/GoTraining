package casting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeAssertion(t *testing.T) {
	var i interface{} = "hello"

	// comma-ok pattern
	s, ok := i.(string)
	assert.True(t, ok)
	assert.Equal(t, "hello", s)

	// failed assertion (safe)
	f, ok := i.(float64)
	assert.False(t, ok)
	assert.Equal(t, 0.0, f) // Returns zero value on failure
}

func TestUnsafeAssertion(t *testing.T) {
	var i interface{} = 42

	// unsafe: PANICS if wrong type
	val := i.(int)
	assert.Equal(t, 42, val)

	// If we did: i.(string) -> panic!
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	_ = i.(string) // This will panic
}

func TestTypeSwitch(t *testing.T) {
	whatIsIt := func(i interface{}) string {
		switch i.(type) {
		case string:
			return "it is a string"
		case int:
			return "it is an int"
		case bool:
			return "it is a bool"
		default:
			return "unknown"
		}
	}

	assert.Equal(t, "it is a string", whatIsIt("hi"))
	assert.Equal(t, "it is an int", whatIsIt(10))
	assert.Equal(t, "it is a bool", whatIsIt(true))
	assert.Equal(t, "unknown", whatIsIt(3.14))
}

type Animal interface {
	Speak() string
}

type Dog struct {
	Breed string
}

func (d Dog) Speak() string { return "Woof!" }

func TestInterfaceAssertion(t *testing.T) {
	// A more practical example: asserting an interface to a concrete struct
	var a Animal = Dog{Breed: "Labrador"}

	// We can speak through the interface
	assert.Equal(t, "Woof!", a.Speak())

	// But if we want to access Breed, we need an assertion
	d, ok := a.(Dog)
	assert.True(t, ok)
	assert.Equal(t, "Labrador", d.Breed)
}
