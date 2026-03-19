# Challenge 06: Package Refactoring

**Difficulty:** Intermediate
**Covers:** `layout`, `init`, `receivers`

---

## Goal

You are given a god-package `util.go` that has accumulated too many responsibilities. Your job is to split it into focused, well-named packages — and replace a hidden `init()` side effect with an explicit setup function.

---

## Starting Point

The `util/` directory contains a single file with:
- `ParseDate(s string) (time.Time, error)` — date parsing
- `Sanitize(s string) string` — string sanitization
- `BuildAuthHeader(token string) string` — HTTP header construction
- A package-level `init()` that sets a global timezone

---

## Requirements

1. **Split into packages** — create appropriately named packages that follow the `layout` module's good-pattern rules:
   - One package per responsibility
   - Package names are nouns, not verbs, not `util`

2. **Replace `init()` with `Setup()`** — the global timezone side effect should become an explicit `Setup(tz string) error` function that callers invoke consciously.

3. **Use pointer receivers where state is mutated** — if your refactored types hold mutable state, use pointer receivers consistently.

4. **Write tests** — prove the refactored code behaves identically to the original:
   - `ParseDate` parses RFC3339 strings correctly
   - `Sanitize` trims whitespace and lowercases
   - `BuildAuthHeader` returns `"Bearer <token>"`
   - `Setup` applies the timezone (test with a known TZ)

---

## The God-Package to Refactor

```go
// util/util.go
package util

import (
    "fmt"
    "net/http"
    "strings"
    "time"
)

var defaultLocation *time.Location

func init() {
    // Side effect: silently sets a global timezone.
    // Any package that imports util gets this side effect.
    defaultLocation, _ = time.LoadLocation("UTC")
}

func ParseDate(s string) (time.Time, error) {
    return time.ParseInLocation(time.RFC3339, s, defaultLocation)
}

func Sanitize(s string) string {
    return strings.ToLower(strings.TrimSpace(s))
}

func BuildAuthHeader(token string) string {
    return fmt.Sprintf("%s %s", http.MethodGet, token) // bug: wrong prefix!
}
```

> **Note:** there is a deliberate bug in `BuildAuthHeader` — fix it as part of the refactor.

---

## Skills Practiced

- Package naming and responsibility separation
- Replacing `init()` with explicit `Setup()` functions
- Pointer vs value receivers
- Writing tests that verify behaviour without coupling to implementation

---

## Hints

- Good package names for this split: `dates`, `sanitize`, `auth` (or `headers`)
- `time.LoadLocation` returns an error — the explicit `Setup` can surface it
- The bug: `BuildAuthHeader` should use `"Bearer"` not `"GET"`
