package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	clientconfig "github.com/romangurevitch/go-training/cmd/client/config"
	"github.com/romangurevitch/go-training/internal/temporal/encryption"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"github.com/romangurevitch/go-training/internal/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	configPath := flag.String("config", "", "path to config file")
	orderPayload := flag.String("order", "", "json order payload")
	workflowName := flag.String("workflow", "ProcessOrder", "workflow to run: ProcessOrder or AutoProcessOrder")
	flag.Parse()

	if *configPath == "" {
		*configPath = "./config/client/local/config.yaml"
	}

	if *orderPayload == "" {
		slog.Error("json order payload is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := clientconfig.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Unable to load config", "error", err)
		os.Exit(1)
	}

	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", cfg.Temporal.Host, cfg.Temporal.Port),
		Logger:   slog.Default(),
		DataConverter: encryption.NewEncryptionDataConverter(
			converter.GetDefaultDataConverter(),
			encryption.DataConverterOptions{Compress: true},
		),
	})
	if err != nil {
		slog.Error("Unable to create client", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	workflowID := "order-" + uuid.New().String()

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: cfg.Temporal.TaskQueueName,
	}

	var o order.Order
	err = json.Unmarshal([]byte(*orderPayload), &o)
	if err != nil {
		slog.Error("Unable to unmarshall payload into order struct", "error", err)
		os.Exit(2)
	}

	var wf interface{}
	switch *workflowName {
	case "AutoProcessOrder":
		wf = workflows.AutoProcessOrder
	default:
		wf = workflows.ProcessOrder
	}

	slog.Info("Starting workflow", "workflow", *workflowName, "workflow_id", workflowID)

	we, err := c.ExecuteWorkflow(context.Background(), options, wf, workflows.Params{
		Order: o,
	})
	if err != nil {
		slog.Error("Unable to execute workflow", "error", err)
		os.Exit(1)
	}

	slog.Info("Workflow started", "workflow_id", we.GetID(), "run_id", we.GetRunID())

	var result order.OrderStatus
	err = we.Get(context.Background(), &result)
	if err != nil {
		slog.Error("Unable get workflow result", "error", err)
		os.Exit(1)
	}

	slog.Info("Workflow completed", "status", result)
}
