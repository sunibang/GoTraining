# Package Layout

Idiomatic Go projects follow a consistent layout that separates concerns.

## Key Directories

| Directory | Purpose |
|---|---|
| `cmd/` | Entry points (one `main` package per binary) |
| `internal/` | Private application code |
| `pkg/` | Public libraries |

## The Layout Problem

Dumping everything in `util` creates a god package with no clear responsibility.

## Better Layout

Name packages by responsibility:
```go
import "myapp/pkg/files"   // files.Open
import "myapp/pkg/strings" // strings.ToUpper
```

## Pitfalls

- Avoid package names like `util`, `common`, `helpers`
- `internal/` enforces boundaries at the Go toolchain level
