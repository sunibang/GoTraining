package pointers

// IncrementValue takes n by value and increments the copy.
// The caller's variable is not affected.
func IncrementValue(n int) int {
	n++
	return n
}

// IncrementPointer takes a pointer to n and increments the original.
func IncrementPointer(n *int) {
	*n++
}

// Counter demonstrates pointer receivers for stateful types.
type Counter struct {
	count int
}

// Increment uses a pointer receiver to mutate the Counter.
func (c *Counter) Increment() {
	c.count++
}

// Value uses a value receiver for a read-only operation.
func (c Counter) Value() int {
	return c.count
}
