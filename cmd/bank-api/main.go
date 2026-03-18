// Package main is the entry point for the Go Bank REST API server.
// See internal/bank/api/ for the HTTP layer implementation.
package main

import (
	"fmt"

	_ "github.com/gin-gonic/gin"
	_ "github.com/go-jet/jet/v2/postgres"
	_ "github.com/lib/pq"
	_ "github.com/spf13/viper"
)

func main() {
	// TODO: initialise config, logger, database store, and start Gin server
	// See internal/bank/api/README.md for implementation guide
	fmt.Println("Go Bank API — not yet implemented")
}
