package challenge

import (
	"math"
	"testing"
)

func TestShapeArea(t *testing.T) {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 10, Height: 5}

	tests := []struct {
		name     string
		shape    Shape
		expected float64
	}{
		{name: "Circle", shape: c, expected: math.Pi * 25},
		{name: "Rectangle", shape: r, expected: 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrintArea(tt.shape); math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("expected area %v, got %v", tt.expected, got)
			}
		})
	}
}
