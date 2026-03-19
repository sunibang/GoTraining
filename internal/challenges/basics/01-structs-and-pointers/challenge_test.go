package challenge

import "testing"

func TestUpdateAge(t *testing.T) {
	p := &Person{Name: "Alice", Age: 25}
	newAge := 26
	UpdateAge(p, newAge)

	if p.Age != newAge {
		t.Errorf("expected age %d, got %d", newAge, p.Age)
	}
}
