// Package calculator provides a simple calculator for testing practice.
package calculator

import "errors"

// Calculator performs basic arithmetic operations.
type Calculator struct{}

// New returns a Calculator. In a real system this might return an error
// (e.g., connecting to a calculation service), which is why we use
// require.NoError in tests rather than just Calculator{}.
func New() (*Calculator, error) {
	return &Calculator{}, nil
}

func (c *Calculator) Add(a, b int) int {
	return a + b
}

func (c *Calculator) Subtract(a, b int) int {
	return a - b
}

// Divide returns a/b or an error when b is zero.
func (c *Calculator) Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
