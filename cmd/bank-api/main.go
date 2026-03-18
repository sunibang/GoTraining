package main

import (
	"log/slog"
	"os"

	"github.com/romangurevitch/go-training/internal/bank/app"
	"github.com/romangurevitch/go-training/internal/bank/config"
)

func init() {
	config.Init()
}

func main() {
	if err := app.Run(); err != nil {
		slog.Error("application failed", slog.Any("error", err))
		os.Exit(1)
	}
}
