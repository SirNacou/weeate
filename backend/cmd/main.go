package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
)

func main() {
	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// Setup Bus
	bus, err := bus.NewBus(logger)
	if err != nil {
		logger.Error("Failed to create bus", slog.String("error", err.Error()))
		os.Exit(1)
	}

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
		bus:    bus,
	}

	// Start bus
	go func() {
		if err := bus.Start(ctx); err != nil {
			logger.Error("Failed to start bus", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Start server
	if err := app.run(app.mount(ctx)); err != nil {
		logger.Error("Failed to run server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
