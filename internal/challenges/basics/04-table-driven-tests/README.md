# Challenge 04: Table-Driven Test Coverage

**Difficulty:** Basic
**Covers:** `testing`, `testify`

---

## Goal

Practice the canonical Go testing pattern: table-driven tests with `t.Run()`, combining `assert` and `require` from testify, and a reusable test helper with `t.Helper()`.

---

## Requirements

A `Calculator` with three methods is provided in `calculator.go`:

| Method | Signature | Error condition |
|--------|-----------|-----------------|
| `Add` | `(a, b int) int` | none |
| `Subtract` | `(a, b int) int` | none |
| `Divide` | `(a, b int) (int, error)` | `b == 0` → error |

Your task in `calculator_test.go`:

1. Write **one table-driven test per method** using `t.Run()`.
2. For `Divide`, include at least one divide-by-zero case.
3. Use `require.NoError` for the setup/connection step (pretend `New()` could fail).
4. Use `assert.Equal` for value comparisons and `assert.Error` for the divide-by-zero case.
5. Extract a helper `checkCalc(t, got, want int)` that uses `t.Helper()` — confirm the failure line in the output points to the test table row, not the helper.

---

## Starter Code

```go
// calculator_test.go
package calculator

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// checkCalc is a reusable assertion helper.
// t.Helper() ensures failure messages point to the caller, not here.
func checkCalc(t *testing.T, got, want int) {
    t.Helper()
    // TODO: assert got == want
}

func TestAdd(t *testing.T) {
    calc, err := New()
    require.NoError(t, err) // stop immediately if construction fails

    tests := []struct {
        name string
        a, b int
        want int
    }{
        // TODO: add test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // TODO: call calc.Add and use checkCalc
        })
    }
}

// TODO: TestSubtract, TestDivide
```

---

## Skills Practiced

- Table-driven tests with `t.Run()`
- `require` vs `assert`: when each matters
- `t.Helper()` for clean failure output

---

## Hints

- `require` stops the test on failure; `assert` continues — use `require` when later assertions are meaningless without an earlier one passing.
- Name your test cases descriptively: `"divide by zero"` is better than `"error case"`.
- Try making a test case fail on purpose and observe where the line number appears with and without `t.Helper()`.
