package challenge

import "errors"

var ErrDivByZero = errors.New("division by zero")

// Divide performs a division of a by b
func Divide(a, b float64) (float64, error) {
	_, _ = a, b
	return 0, nil // TODO: implement this
}
