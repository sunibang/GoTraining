// Package main is the entry point for the Temporal worker.
// See internal/temporal/README.md for context.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	workerconfig "github.com/romangurevitch/go-training/cmd/worker/config"
	"github.com/romangurevitch/go-training/internal/temporal/activities"
	"github.com/romangurevitch/go-training/internal/temporal/encryption"
	"github.com/romangurevitch/go-training/internal/temporal/integrations/inventory"
	"github.com/romangurevitch/go-training/internal/temporal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		*configPath = "./config/worker/local/config.yaml"
	}

	cfg, err := workerconfig.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Unable to load config", "error", err)
		os.Exit(1)
	}

	// Create the Temporal client,
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", cfg.Temporal.Host, cfg.Temporal.Port),
		Logger:   slog.Default(),
		// Set DataConverter to ensure that workflow inputs and results are
		// encrypted/decrypted as required.
		DataConverter: encryption.NewEncryptionDataConverter(
			converter.GetDefaultDataConverter(),
			encryption.DataConverterOptions{Compress: true},
		),
		// Use a ContextPropagator so that the KeyID value set in the workflow context is
		// also availble in the context for activities.
		ContextPropagators: []workflow.ContextPropagator{encryption.NewContextPropagator()},
	})
	if err != nil {
		slog.Error("Unable to create Temporal client", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// Create the Temporal worker,
	w := worker.New(c, cfg.Temporal.TaskQueueName, worker.Options{})

	// inject HTTP client into the Activities Struct,
	acts := activities.NewOrderActivities(inventory.NewClient(cfg.InventoryAPI.BaseURL))

	// Register Workflow and Activities
	w.RegisterWorkflow(workflows.ProcessOrder)
	w.RegisterWorkflow(workflows.ProcessPayment)
	w.RegisterWorkflow(workflows.AutoProcessOrder)
	w.RegisterActivity(acts)

	// Start the Worker
	if err := w.Run(worker.InterruptCh()); err != nil {
		slog.Default().Error("Unable to start Temporal worker", "error", err)
		os.Exit(1)
	}

	slog.Info("Shutting down...")
}
