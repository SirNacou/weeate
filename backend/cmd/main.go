package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/data"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/sync/errgroup"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	slog.SetDefault(logger)

	// Setup configuration
	env, err := configs.LoadEnv()
	if err != nil {
		logger.Error("Failed to load environment variables", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Setup context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Database connection
	db, err := data.ConnectToPostgres(ctx, env.GetDBDsn())
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := data.MigratePostgresDB(db, ); err != nil {
		slog.Error("Failed to migrate database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("Failed to get sql.DB")
		os.Exit(1)
	}

	// conn, err := sqlDB.Conn(ctx)
	// if err != nil {
	// 	slog.Error("Failed to connect to pgx", slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }
	// defer conn.Close()

	// Setup Supabase auth
	supabaseClient, err := supabase.NewClient(env.SUPABASE_URL, env.SUPABASE_API_KEY, &supabase.ClientOptions{})
	if err != nil {
		slog.Error("Failed to initalize the client: ", slog.String("error", err.Error()))
	}
	supabaseService := auth.NewSupabaseService(supabaseClient)

	// Setup Bus
	bus, err := bus.NewBus(sqlDB, logger)
	if err != nil {
		logger.Error("Failed to create bus", slog.String("error", err.Error()))
		os.Exit(1)
	}

	server, err := newServer(ctx, env, logger, db, supabaseService, bus)
	if err != nil {
		slog.Error("Failed to create server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	g, gCtx := errgroup.WithContext(ctx)

	// Start Bus
	g.Go(func() error {
		return bus.Start(gCtx)
	})

	// Start HTTP server
	g.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%v", env.PORT), server)
	})

	// Start server
	if err := g.Wait(); err != nil {
		logger.Error("Failed to run server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
