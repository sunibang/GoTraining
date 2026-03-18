package main

import (
	"net/http"
	"os"

	"github.com/romangurevitch/go-training/internal/bank/cli"
	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func main() {
	// API base URL can be configured via environment variable for flexibility.
	apiURL := os.Getenv("BANK_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080" // Root URL, client handles /v1/
	}

	// Initialize the Bank REST client.
	bankClient := bank.New(apiURL, &http.Client{})

	// Load token from environment if present.
	if token := os.Getenv("BANK_TOKEN"); token != "" {
		bankClient.SetToken(token)
	}

	// Execute the CLI.
	cobra.CheckErr(cli.New(bankClient).Execute())
}
