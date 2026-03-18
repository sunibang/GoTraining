package account

import (
	"fmt"

	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func getBalanceCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "balance [account-id]",
		Short: "Check the balance of an account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// QUEST 5: Participants should implement this method.
			// It should call bankClient.GetAccount and print the balance.
			fmt.Println("Balance check not implemented yet. Implement it in Quest 5!")
		},
	}
}
