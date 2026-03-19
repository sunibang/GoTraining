# 🧬 Generics in Go

Generics (introduced in Go 1.18) allow you to write code that works with multiple types while maintaining complete type safety. This eliminates the need for repetitive code or unsafe `interface{}`-based casting.

## 📌 Why use Generics?
- **DRY (Don't Repeat Yourself)**: One function can work for `int`, `float`, and `string`.
- **Type Safety**: The compiler catches errors early, so you don't need to check types at runtime.
- **Cleaner APIs**: No more `map[string]interface{}` or manual type assertions.
- **Performance**: Faster than using the empty interface because there's no runtime reflection overhead.

---

## 🏗️ How Generics Work (The Blueprint)

Generics act as a "blueprint". The compiler generates a specific version of your code for each type you use (Monomorphization).

```text
  +-------------------------------------------------------+
  |              Generic Blueprint (Source)               |
  |      func Print[T any](v T) { ... }                   |
  +-------------------------------------------------------+
            |
            v
  +-------------------------------------------------------+
  |              Compiler Magic (Output)                  |
  |  [Print(42)]   -> func PrintInt(v int)                |
  |  [Print("Hi")] -> func PrintStr(v string)             |
  +-------------------------------------------------------+
```

---

## ✍️ Anatomy of a Generic Function

Generics use square brackets `[]` to define type parameters and constraints.

```go
// 1. Type Parameter [T any]
// 2. Type Constraint (any, comparable, Ordered)
func Min[T comparable](a, b T) T {
    if a == b { // Comparable allows == and !=
        return a
    }
    return b
}

// 3. Custom Constraint (Unions)
type Number interface {
    ~int | ~float64 // ~ includes underlying types
}
```

---

## 🏃 Common Operations & Patterns

| Category | Pattern | Example |
|---------|---------|---------|
| **Functions** | Utility Logic | `func Sum[T Number](vals []T) T` |
| **Structs** | Containers | `type Stack[T any] struct { ... }` |
| **Functional**| Data Processing | `Filter`, `Reduce`, `Map` |
| **Sets** | Deduplication | `type Set[T comparable] struct { ... }` |

---

## 🚀 Advanced Patterns (Included in this folder)

### 1. The Zero Value Pattern
When returning from a generic function that might "fail" (like popping an empty stack), you can't return `nil` for all types. Use `var zero T`:
```go
func First[T any](s []T) (T, bool) {
    if len(s) == 0 {
        var zero T // Returns 0, "", false, or nil depending on T
        return zero, false
    }
    return s[0], true
}
```

### 2. Generic Interfaces
Interfaces can also be parameterized. This describes **behaviour on T** rather than restricting what T is:
```go
type Container[T any] interface {
    Add(T)
    Get() T
}
```

### 3. Multiple Type Parameters
Types like `Pair[T, U]` or functions like `Reduce[T, U]` allow different types to interact safely:
```go
// Reducing a slice of ints into a single string
result := Reduce([]int{1, 2}, "Count: ", func(acc string, val int) string {
    return acc + strconv.Itoa(val)
})
```

---

## 💡 Pro Tips for Starters

### 1. The `~` Tilde is Your Friend
Always use `~` in your constraints (e.g., `~int` instead of `int`) if you want to support custom types that are aliases of those built-ins.

### 2. Type Inference
The compiler is smart! You don't always need to specify the type in brackets.
- `Sum[int]([]int{1, 2})` (Explicit)
- `Sum([]int{1, 2})` (Inferred - **Preferred**)

### 3. When NOT to use Generics
If your logic only works for one or two types, or if you find yourself using `any` and then type-asserting inside the function, you probably don't need generics.

---

## 🛠️ Practical Examples in this Directory

- `generics.go`: Basic `Stack` and `SliceContains`.
- `advanced_types.go`: `Pair`, `Set`, `Container` interface, and `ZeroValue` pattern.
- `functional.go`: `Filter` and `Reduce` implementations.
- `constraints.go`: Custom `Number` and `Ordered` constraints.

**Run the tests to see it all in action!**
```bash
go test -v ./internal/basics/generics/...
```
