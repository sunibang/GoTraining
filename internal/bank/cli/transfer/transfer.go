package transfer

import (
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

// GetTransferCmd returns the 'transfer' command group.
func GetTransferCmd(bankClient bank.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer funds between accounts",
	}

	cmd.AddCommand(getCreateTransferCmd(bankClient))

	return cmd
}

func getCreateTransferCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "create [from-id] [to-id] [amount]",
		Short: "Create a new transfer",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the CLI logic to execute a transfer.
			// 1. Parse arguments (fromID, toID, amount)
			// 2. Call bankClient.Transfer
			// 3. Print the success response as JSON
		},
	}
}
