package cli

import (
	"github.com/romangurevitch/go-training/internal/bank/cli/account"
	"github.com/romangurevitch/go-training/internal/bank/cli/transfer"
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

// New creates and assembles the bank-cli root command.
func New(bankClient bank.Client) *cobra.Command {
	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "bank-cli",
		Short: "Go Bank CLI",
		Long: `A well-structured example of a command-line interface 
for interacting with the Go Bank API.`,
	}

	rootCmd.AddCommand(account.GetAccountCmd(bankClient))
	rootCmd.AddCommand(transfer.GetTransferCmd(bankClient))

	return rootCmd
}
