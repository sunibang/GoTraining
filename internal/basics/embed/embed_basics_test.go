package embed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1. Basic Struct Embedding (Composition)
// The fields and methods of the embedded struct are "promoted" to the outer struct.
type User struct {
	Name string
}

func (u User) Greet() string {
	return "Hi, I'm " + u.Name
}

type Admin struct {
	User  // <--- Embedding! No field name.
	Level int
}

func TestStructEmbedding_Promotion(t *testing.T) {
	a := Admin{
		User:  User{Name: "Alice"},
		Level: 10,
	}

	// 1. Promoted Fields: Can access Name directly on Admin
	assert.Equal(t, "Alice", a.Name)

	// 2. Can still access through the inner name (the type name)
	assert.Equal(t, "Alice", a.User.Name) //nolint:staticcheck // intentionally showing explicit embedded field access

	// 3. Promoted Methods: Can call Greet() directly on Admin
	assert.Equal(t, "Hi, I'm Alice", a.Greet())
}

// 2. Promotion to Interfaces
// If the embedded type satisfies an interface, the outer type also does.
type Greeter interface {
	Greet() string
}

func TestStructEmbedding_Interfaces(t *testing.T) {
	a := Admin{User: User{Name: "Bob"}}

	// Admin satisfies Greeter because User does
	var g Greeter = a
	assert.Equal(t, "Hi, I'm Bob", g.Greet())
}
