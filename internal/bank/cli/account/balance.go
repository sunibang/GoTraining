package account

import (
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func getBalanceCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "balance [account-id]",
		Short: "Check the balance of an account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement the CLI logic to check an account balance.
			// 1. Parse arguments (account-id)
			// 2. Call bankClient.GetAccount
			// 3. Print the success balance response as json
		},
	}
}
