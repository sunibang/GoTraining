package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkCalc is a reusable assertion helper.
// t.Helper() ensures failure messages point to the caller's line, not here.
func checkCalc(t *testing.T, got, want int) {
	t.Helper()
	// TODO: assert got == want using assert.Equal
	assert.Equal(t, want, got)
}

func TestAdd(t *testing.T) {
	calc, err := New()
	require.NoError(t, err, "calculator construction must not fail")

	tests := []struct {
		name string
		a, b int
		want int
	}{
		// TODO: add at least 3 test cases (positive, zero, negative)
		{"TODO: positive numbers", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Add(tt.a, tt.b)
			checkCalc(t, got, tt.want)
		})
	}
}

func TestSubtract(t *testing.T) {
	calc, err := New()
	require.NoError(t, err)

	tests := []struct {
		name string
		a, b int
		want int
	}{
		// TODO: add test cases
		{"TODO: implement me", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calc.Subtract(tt.a, tt.b)
			checkCalc(t, got, tt.want)
		})
	}
}

func TestDivide(t *testing.T) {
	calc, err := New()
	require.NoError(t, err)

	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		// TODO: add success cases AND divide-by-zero case
		{"TODO: implement me", 0, 1, 0, false},
		{"TODO: divide by zero", 5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calc.Divide(tt.a, tt.b)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			checkCalc(t, got, tt.want)
		})
	}
}
