package account

import (
	"context"
	"fmt"

	"github.com/romangurevitch/go-training/internal/pkg/json"
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func getCreateCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "create [owner]",
		Short: "Create a new bank account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			owner := args[0]
			acc, err := bankClient.CreateAccount(context.Background(), owner)
			cobra.CheckErr(err)
			fmt.Println(json.ToJSONString(acc))
		},
	}
}
