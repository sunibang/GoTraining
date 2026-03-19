package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// Create a new user
	u := New("1", "Alice", "ALICE@example.com")

	// 1. Verify basic data promotion (initialization)
	assert.Equal(t, "1", u.GetID())
	assert.Equal(t, "Alice", u.GetName())
	assert.Equal(t, "alice@example.com", u.GetEmail()) // Verify it was lowercased

	// 2. Verify default role (not admin)
	assert.False(t, u.IsAdmin())

	// 3. Verify Stringer implementation
	assert.Equal(t, "User[1]: Alice <alice@example.com>", u.String())
}

func TestEncapsulation(t *testing.T) {
	u := New("2", "Bob", "bob@example.com")

	// We can't access 'role' directly on the interface
	// u.role would be a compile error!

	// But we can interact through authorized methods
	concreteUser, ok := u.(*user)
	assert.True(t, ok)

	// Internally we can promote
	concreteUser.PromoteToAdmin()
	assert.True(t, u.IsAdmin())
}

func TestStructTags(t *testing.T) {
	// This test demonstrates how struct tags are used by the json package.
	u := &user{ID: "3", Name: "Charlie", Email: "charlie@example.com"}
	b, err := json.Marshal(u)
	assert.NoError(t, err)
	// The unexported 'role' field is not included in the JSON output,
	// and the exported fields use the names from the `json` tags.
	expectedJSON := `{"id":"3","name":"Charlie","email":"charlie@example.com"}`
	assert.JSONEq(t, expectedJSON, string(b))
}
