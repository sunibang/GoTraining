# 🚀 Entities & Structs in Go

In Go, an **Entity** is a `struct` that groups related data with **methods** that define its behaviour. Unlike traditional OOP, Go uses composition and visibility rules instead of class inheritance.

## 📌 Why use Entities?
- **Data Grouping**: Keep related fields together in a single type.
- **Encapsulation**: Control what is public (Exported) or private (Unexported).
- **Consistency**: Use factory functions (`New`) to ensure objects are always valid.
- **Behaviour**: Attach logic directly to data using methods.

---

## 🏗️ How Entities Work

Go combines data and logic without the complexity of classes.

```text
  +-------------------------------------------------------+
  |                   Entity Structure                    |
  +-------------------------------------------------------+
  |  1. Data (Struct): Fields like ID, Name, Email        |
  |  2. Factory (New): Function to safely create object   |
  |  3. Behaviour (Methods): Logic that operates on data  |
  +-------------------------------------------------------+
            |
            v
  +-------------------------------------------------------+
  | Result: A self-contained, testable domain object      |
  +-------------------------------------------------------+
```

---

## ✍️ Anatomy of an Entity

Entities typically follow these conventions:
1. Defined as a `struct`.
2. Use capitalization for visibility (Upper = Public, Lower = Private).
3. Have a `New` factory function.

```go
// 1. The Struct (Data)
type User struct {
    ID    int    // Public
    Email string // Public
    role  string // Private (internal to package)
}

// 2. The Factory (Constructor)
func NewUser(id int, email string) *User {
    return &User{
        ID:    id,
        Email: email,
        role:  "guest", // Default value
    }
}

// 3. The Method (Behaviour)
func (u *User) IsAdmin() bool {
    return u.role == "admin"
}
```

---

## 🏃 Common Operations

| Task | Pattern | Example |
|---------|-------------|---------|
| **Creation** | Factory Function | `u := NewUser(1, "test@example.com")` |
| **JSON** | Struct Tags | ``ID int `json:"id"` `` |
| **Updates** | Pointer Receiver | `func (u *User) SetEmail(e string)` |
| **Read-only** | Value Receiver | `func (u User) GetID() int` |

---

## 💡 Pro Tips for Starters

### 1. Pointer vs. Value Receivers
- Use **Pointer Receivers** (`*User`) if you need to modify the struct or if it's large.
- Use **Value Receivers** (`User`) for small, immutable types or when no modification is needed.

### 2. The "New" Pattern
Always return an interface or a pointer from your factory function to hide implementation details and allow for easier mocking in tests.

### 3. Struct Tags for APIs
Use struct tags to map Go fields to external formats like JSON or Database columns.
```go
type Account struct {
    Balance float64 `json:"balance" db:"total_amount"`
}

// New is the factory function
func New(name string) User {
    return &user{
        Name: name,
        role: "customer", // Default
    }
}

func (u *user) GetName() string { return u.Name }
func (u *user) IsAdmin() bool   { return u.role == "admin" }
```

---

## 🛠️ Practical Example: User Entity

In `entity.go`, we demonstrate a complete User entity with:
- Private fields for security.
- A public interface for abstraction.
- Methods for business logic.

**Try running the tests to see it in action!**
```bash
go test -v ./internal/basics/entity/...
```

---

## 📚 Reference

| Feature | Go Keyword/Convention |
|---------|-----------------------|
| **Public** | `CapitalizedName` |
| **Private** | `lowercaseName` |
| **Methods** | `func (receiver) Name()` |
| **Composition** | Struct Embedding |
