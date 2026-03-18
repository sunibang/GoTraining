// Package main is the entry point for the Go Bank CLI tool.
// See internal/bank/ for domain models and service layer.
package main

import (
	"fmt"

	_ "github.com/spf13/cobra"
	_ "golang.org/x/sync/errgroup"
)

func main() {
	// TODO: wire up Cobra root command
	// See internal/bank/README.md for implementation guide
	fmt.Println("Go Bank CLI — not yet implemented")
}
