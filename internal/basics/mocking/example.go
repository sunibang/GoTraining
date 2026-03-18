package mocking

import "github.com/romangurevitch/go-training/internal/basics/mocking/calculator"

func ExampleFunction(adder calculator.Adder, x int, y int) (int, error) {
	result, err := adder.SingleDigitAdd(x, y)
	return result, err
}
