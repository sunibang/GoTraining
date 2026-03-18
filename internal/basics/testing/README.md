# Testing

Go's `testing` package provides everything needed for unit, integration, and benchmark tests.

## Table-Driven Tests

The standard Go pattern for testing multiple cases:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, Add(tt.a, tt.b))
        })
    }
}
```

## Running Tests

```bash
go test ./...                 # all tests
go test -v ./...              # verbose
go test -race ./...           # detect race conditions
go test -cover ./...          # coverage
```

## Pitfalls

- Name test cases clearly
- Use `t.Helper()` in helper functions
- Use subtests (`t.Run`) to isolate failures
