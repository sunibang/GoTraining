package challenge

import (
	"errors"
	"math"
	"testing"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
		err      error
	}{
		{name: "success", a: 10, b: 2, expected: 5, err: nil},
		{name: "fractional result", a: 1, b: 3, expected: 1.0 / 3.0, err: nil},
		{name: "division by zero", a: 10, b: 0, expected: 0, err: ErrDivByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)

			if tt.err != nil {
				if !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v, got %v", tt.err, err)
				}
				// Optional: checking the returned value on error is often not required
				// but here we ensure it matches the expectation if specified.
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("expected result %v, got %v", tt.expected, got)
			}
		})
	}
}
