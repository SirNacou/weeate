package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
)

func main() {
	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// Setup configuration
	env, err := configs.LoadEnv()
	if err != nil {
		logger.Error("Failed to load environment variables", slog.String("error", err.Error()))
		os.Exit(1)
	}

	config := newConfig(env)
	app := application{
		config: config,
		logger: logger,
	}

	// Start server
	if err := app.run(app.mount(ctx)); err != nil {
		logger.Error("Failed to run server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
