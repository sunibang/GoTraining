package generics

import (
	"fmt"
	"strings"
)

// Number is a custom constraint using a Union type.
// The '~' symbol means "any type whose underlying type is int or float64".
// This is crucial for supporting custom types like 'type MyInt int'.
type Number interface {
	~int | ~int64 | ~float64
}

// Sum adds all elements in a slice of Numbers.
func Sum[T Number](values []T) T {
	var total T
	for _, v := range values {
		total += v
	}
	return total
}

// Map converts a slice of one type to a slice of another type.
func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Stringer is a basic interface constraint.
type Stringer interface {
	String() string
}

// Magic is a struct for demonstration.
type Magic struct {
	Name  string
	Spell []string
}

func (m *Magic) String() string {
	return fmt.Sprintf("This is a %s. Repeat after me: %s", m.Name, strings.Join(m.Spell, " "))
}

// CustomError is a struct for demonstration.
type CustomError struct {
	HTTPStatusCode int
	ErrorMessage   string
}

func (e *CustomError) String() string {
	return fmt.Sprintf("Error: status code %d, error message %s", e.HTTPStatusCode, e.ErrorMessage)
}

func ToString[T Stringer](val T) {
	fmt.Println(val.String())
}
