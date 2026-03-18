package account

import (
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

// GetAccountCmd returns the 'account' command group.
func GetAccountCmd(bankClient bank.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Account related commands",
	}

	cmd.AddCommand(getCreateCmd(bankClient))
	cmd.AddCommand(getBalanceCmd(bankClient))

	return cmd
}
