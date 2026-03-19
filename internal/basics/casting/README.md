# 🔄 Type Conversion & Assertion in Go

In Go, there are two primary ways to change or interpret types: **Type Conversion** (between concrete types) and **Type Assertion** (from interfaces to concrete types).

---

## 1. Type Conversion (Concrete to Concrete)

Type conversion is used when you have two compatible concrete types (like `int32` and `int64`) and you want to convert a value from one to the other.

### 🖼️ Pictorial Representation
```text
  +-------------+                     +-------------+
  |   Source    |     T(value)        |    Target   |
  |    Type     |  -------------->    |     Type    |
  | (e.g. int)  |                     | (e.g. float)|
  +-------------+                     +-------------+
         |                                   |
         v                                   v
       [ 42 ]           ------>           [ 42.0 ]
```

### 📝 Example
```go
var i int = 42
var f float64 = float64(i) // Manual conversion required
```

> [!WARNING]
> Go **never** performs implicit type conversion. Even `int` and `int64` are different types and require manual conversion.

---

## 2. Type Assertion (Interface to Concrete)

Type assertion is used to extract a concrete value from an interface. It "asserts" that the interface holds a specific type.

### 🖼️ Pictorial Representation
```text
  +-----------------------+
  |      Interface        |
  |  (Dynamic Type + Val) |
  +-----------------------+
              |
              |  v, ok := i.(T)
              v
     /-----------------\
    |   Is it Type T?   |
     \-----------------/
      /               \
    [Yes]            [No]
      |                |
      v                v
 v = Value        v = Zero Value
 ok = true        ok = false
```

### 📝 Example
```go
var i interface{} = "hello"

// Safe assertion
s, ok := i.(string) 

// Unsafe assertion (PANICS if not a string!)
s := i.(string) 
```

---

## 3. Type Switches

A type switch is a cleaner way to handle multiple possible types for an interface.

```go
switch v := i.(type) {
case string:
    fmt.Printf("It's a string: %s\n", v)
case int:
    fmt.Printf("It's an int: %d\n", v)
default:
    fmt.Printf("Unknown type!\n")
}
```

---

## 🧪 Running the Examples

You can find runnable examples in the test files:
- `conversion_test.go`: Examples of numeric and string conversions.
- `assertion_test.go`: Examples of interface assertions and type switches.

```bash
go test -v ./internal/basics/casting/...
```
