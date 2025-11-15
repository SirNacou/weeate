package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/SirNacou/weeate/backend/internal/api/auth"
	"github.com/SirNacou/weeate/backend/internal/api/common"
	"github.com/SirNacou/weeate/backend/internal/api/foods"
	domain "github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/db"
	"github.com/SirNacou/weeate/backend/internal/usecase"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/supabase-community/supabase-go"
)

type application struct {
	config config
	logger *slog.Logger
}

func (a *application) mount(ctx context.Context) http.Handler {
	// Setup Supabase auth
	supabaseClient, err := supabase.NewClient(a.config.SUPABASE_URL, a.config.SUPABASE_API_KEY, &supabase.ClientOptions{})
	if err != nil {
		slog.Error("Failed to initalize the client: ", slog.String("error", err.Error()))
	}

	// Database connection
	db, err := db.ConnectToPostgres(ctx, a.config.DSN)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db.AutoMigrate(&domain.Food{})

	// Setup application handlers
	handlers := usecase.NewHandlers(db, supabaseClient)

	// Setup Fiber app
	app := fiber.New(fiber.Config{})

	app.Use(slogfiber.New(a.logger))

	app.Use(recover.New())

	app.Use(common.CORSMiddleware(a.config.GO_ENV))
	authMiddleware, err := common.AuthMiddleware(ctx, a.config.SUPABASE_AUTH_URL, a.config.SUPABASE_COOKIE_AUTH_NAME)
	if err != nil {
		slog.Error("Failed to initialize auth middleware", slog.String("error", err.Error()))
		os.Exit(1)
	}
	app.Use(authMiddleware)

	api := humafiber.New(app, huma.DefaultConfig("Weeate API", "v1.0.0"))
	foodsEndpoint := foods.NewFoodsEndpoint(handlers)
	foodsEndpoint.Register(api)

	huma.Get(api, "/", func(ctx context.Context, i *struct{}) (*auth.User, error) {
		user, err := auth.GetUserContext(ctx)
		if err != nil {
			return nil, huma.Error401Unauthorized("Unauthorized", err)
		}
		return &user, nil
	})

	return api.Adapter()
}

func (a *application) run(h http.Handler) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", a.config.PORT), h)
}

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

type config struct {
	PORT                      int
	Timezone                  string
	DSN                       string
	SUPABASE_URL              string
	SUPABASE_AUTH_URL         string
	SUPABASE_API_KEY          string
	SUPABASE_COOKIE_AUTH_NAME string
	GO_ENV                    string
	IMAGE_KIT_API_KEY         string
	IMAGEKIT_URL              string
}

func newConfig(e configs.Env) config {
	return config{
		PORT:                      e.PORT,
		Timezone:                  e.Timezone,
		DSN:                       e.GetDBDsn(),
		SUPABASE_URL:              e.SUPABASE_URL,
		SUPABASE_AUTH_URL:         e.SUPABASE_AUTH_URL,
		SUPABASE_API_KEY:          e.SUPABASE_API_KEY,
		SUPABASE_COOKIE_AUTH_NAME: e.SUPABASE_COOKIE_AUTH_NAME,
		GO_ENV:                    e.GO_ENV,
		IMAGE_KIT_API_KEY:         e.IMAGE_KIT_API_KEY,
		IMAGEKIT_URL:              e.IMAGEKIT_URL,
	}
}
