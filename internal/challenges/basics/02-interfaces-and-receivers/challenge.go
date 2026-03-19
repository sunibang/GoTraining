package challenge

import "math"

// Shape represents a geometric shape
type Shape interface {
	Area() float64
}

// Circle represents a circle
type Circle struct {
	Radius float64
}

// Area should return the area of the circle
func (c Circle) Area() float64 {
	return 0 // TODO: implement this
}

// Rectangle represents a rectangle
type Rectangle struct {
	Width, Height float64
}

// Area should return the area of the rectangle
func (r Rectangle) Area() float64 {
	return 0 // TODO: implement this
}

// PrintArea returns the area of a shape
func PrintArea(s Shape) float64 {
	_ = math.Pi     // Use math to avoid unused import error
	return s.Area() // Call the method
}
