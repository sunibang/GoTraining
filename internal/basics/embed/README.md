# 📦 Embedding in Go

Go uses the term "embedding" in two different contexts: **Struct Embedding** (for composition) and the **`go:embed` directive** (for static assets). Both are powerful features that simplify your code and deployments.

---

## 1. Struct Embedding (Composition)

Go does not have classes or inheritance. Instead, it uses **composition** through struct embedding. When you embed one struct into another, the outer struct "borrows" the fields and methods of the inner one.

### 🖼️ Conceptual View
```text
  Traditional Inheritance        Go Struct Embedding
  +-----------+                  +-----------+
  |   Base    |                  |   Outer   |
  +-----------+                  | +-------+ |
        ^                        | | Inner | |
        | (Is-a)                 | +-------+ | (Has-a, looks like Is-a)
  +-----------+                  +-----------+
  |   Child   |
  +-----------+
```

### 📝 Example
```go
type User struct {
    Name string
}
func (u User) Greet() string { return "Hi, I'm " + u.Name }

type Admin struct {
    User  // <--- Embedded! No field name.
    Level int
}

a := Admin{User: User{Name: "Alice"}, Level: 1}
fmt.Println(a.Greet()) // "Hi, I'm Alice" (Promoted method)
fmt.Println(a.Name)    // "Alice" (Promoted field)
```

---

## 2. File Embedding (`go:embed`)

Introduced in Go 1.16, the `embed` package allows you to include static files (like config, HTML, or images) directly into your compiled binary. No more missing config files in production!

### 🖼️ The Process
```text
  Source Code + Asset Files         Single Binary
  +-------------+                   +-------------+
  | main.go     |                   |  [Binary]   |
  | config.json |  --- go build --> |  [Data]     |
  | logo.png    |                   |  [Logic]    |
  +-------------+                   +-------------+
                                     (Runs anywhere!)
```

### 📝 Example
```go
import "embed"

//go:embed hello.txt
var hello string

//go:embed config/*.yaml
var configFS embed.FS // Embed multiple files as a filesystem
```

---

## 3. Key Differences & Pitfalls

### Struct Embedding
- **Shadowing**: If both structs have a field named `ID`, the outer one "wins." You can still access the inner one via `outer.Inner.ID`.
- **Not Subtyping**: An `Admin` is **not** a `User`. You cannot pass an `Admin` to a function expecting a `User`.

### File Embedding
- **Read-Only**: Files embedded via `go:embed` are read-only at runtime.
- **Global Only**: The `//go:embed` directive only works on global (package-level) variables.
- **Import Required**: You must `import _ "embed"` (or just `import "embed"`) even if you don't use the `embed.FS` type.

---

## 🧪 Running the Examples

Explore the unit tests for runnable patterns:
- `embed_basics_test.go`: Struct and method promotion.
- `embed_playground_test.go`: Shadowing and overshadowing methods.
- `embed_file_test.go`: Using the `go:embed` directive.

```bash
go test -v ./internal/basics/embed/...
```
