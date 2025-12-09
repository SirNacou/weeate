package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/data"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/imagekit"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// Setup configuration
	config, err := configs.LoadEnv()
	if err != nil {
		logger.Error("Failed to load environment variables", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Setup context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Database connection
	db, err := data.ConnectToPostgres(ctx, config.GetDBDsn())
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := data.MigratePostgresDB(db); err != nil {
		slog.Error("Failed to migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB")
		os.Exit(1)
	}

	// Setup Supabase auth
	supabaseService, err := auth.NewSupabaseService(config.SUPABASE_URL, config.SUPABASE_API_KEY)
	if err != nil {
		slog.Error("Failed to initalize the client: ", slog.String("error", err.Error()))
		os.Exit(1)
	}

	imagekitClient := imagekit.NewImageKitClient(config.IMAGE_KIT_API_KEY, config.IMAGEKIT_WEBHOOK_KEY)

	// Setup Bus
	bus, err := bus.NewBus(sqlDB, logger)
	if err != nil {
		logger.Error("Failed to create bus", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Setup Centrifugo gRPC client
	centrifugoClient, err := centrifugo.NewCentrifugoClient(config)
	if err != nil {
		slog.Error("Failed to create Centrifugo gRPC client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer centrifugoClient.Close()

	// Run server
	srv := NewServer(config, db, supabaseService, bus, centrifugoClient, imagekitClient)
	if err := srv.Run(ctx); err != nil {
		logger.Error("Failed to run server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
