package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods"
	"github.com/SirNacou/weeate/backend/internal/features/orders"
	"github.com/SirNacou/weeate/backend/internal/features/polls"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"gorm.io/gorm"
)

func newServer(ctx context.Context, config configs.Env, logger *slog.Logger, db *gorm.DB, supabaseService *auth.SupabaseService, bus *bus.Bus) (http.Handler, error) {
	// Setup Fiber app
	app := fiber.New(fiber.Config{})

	app.Use(slogfiber.New(logger))
	app.Use(recover.New())
	app.Use(api.CORSMiddleware(config.GO_ENV))
	authMiddleware, err := api.AuthMiddleware(ctx, config.SUPABASE_AUTH_URL, config.SUPABASE_COOKIE_AUTH_NAME)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize auth middleware: %v", err)
	}
	app.Use(authMiddleware)

	api := humafiber.New(app, huma.DefaultConfig("Weeate API", "v1.0.0"))

	foods.RegisterFoodsModule(api, db, supabaseService)
	polls.RegisterPollsModule(api, bus, db, supabaseService)
	orders.RegisterOrdersModule(api, bus, db, supabaseService)

	huma.Get(api, "/", func(ctx context.Context, i *struct{}) (*auth.User, error) {
		user, err := auth.GetUserContext(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized", err)
		}
		return &user, nil
	})

	return api.Adapter(), nil
}
