package calculator

import "errors"

// generate mocks using mockgen
//go:generate mockgen -destination=mocks/mock_adder.go -package=mocks github.com/romangurevitch/go-training/internal/basics/mocking/calculator Adder

type Adder interface {
	SingleDigitAdd(x int, y int) (int, error)
}

type calculator struct{}

func (c *calculator) SingleDigitAdd(x int, y int) (int, error) {
	if x >= 10 || y >= 10 {
		return 0, errors.New("too high")
	}
	return x + y, nil
}
