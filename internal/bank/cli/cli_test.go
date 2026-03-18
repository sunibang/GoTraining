package cli

import (
	"testing"

	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/stretchr/testify/assert"
)

// mockClient is a simple manual mock for testing command registration.
type mockClient struct {
	bank.Client
}

func TestNewCLI(t *testing.T) {
	m := &mockClient{}
	rootCmd := New(m)

	assert.Equal(t, "bank-cli", rootCmd.Use)
	assert.True(t, rootCmd.HasSubCommands())

	// Verify 'account' command is registered
	accCmd, _, err := rootCmd.Find([]string{"account"})
	assert.NoError(t, err)
	assert.Equal(t, "account", accCmd.Use)

	// Verify 'account create' command is registered
	createCmd, _, err := rootCmd.Find([]string{"account", "create"})
	assert.NoError(t, err)
	assert.Equal(t, "create [owner]", createCmd.Use)

	// Verify 'account balance' command is registered
	balanceCmd, _, err := rootCmd.Find([]string{"account", "balance"})
	assert.NoError(t, err)
	assert.Equal(t, "balance [account-id]", balanceCmd.Use)
}
