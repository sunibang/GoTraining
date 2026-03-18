package calculator

import "errors"

type Adder interface {
	SingleDigitAdd(x int, y int) (int, error)
}

type calculator struct{}

func New() Adder {
	return &calculator{}
}

func (c *calculator) SingleDigitAdd(x int, y int) (int, error) {
	if x >= 10 || y >= 10 {
		return 0, errors.New("too high")
	}
	return x + y, nil
}
