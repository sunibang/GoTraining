# Pointers

A pointer holds the memory address of a value. Go passes everything by value — pointers are how you share and mutate data across function boundaries.

## Value vs Pointer

```go
func incrementValue(n int) {
    n++           // modifies a copy — caller sees no change
}

func incrementPointer(n *int) {
    *n++          // dereferences and modifies the original
}

x := 10
incrementValue(x)    // x is still 10
incrementPointer(&x) // x is now 11
```

## When to Use Pointers

- When a function needs to **mutate** the caller's variable
- When copying a large struct is expensive
- When you need to express **optional** values (a nil pointer means "absent")

## Pointer Receivers vs Value Receivers

```go
type Counter struct{ count int }

func (c *Counter) Increment() { c.count++ } // pointer receiver — mutates c
func (c Counter) Value() int  { return c.count } // value receiver — read-only
```

Rule of thumb: if any method needs a pointer receiver, make **all** methods pointer receivers for consistency.

## Pitfalls

- Dereferencing a nil pointer panics: `var p *int; *p = 1` → panic
- Returning a pointer to a local variable is safe in Go (unlike C) — the GC keeps it alive
- Don't over-use pointers for small types (`int`, `bool`, `string`) — the overhead isn't worth it
