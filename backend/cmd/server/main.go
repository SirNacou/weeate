package main

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"

	"github.com/SirNacou/weeate/backend/internal/api/auth"
	"github.com/SirNacou/weeate/backend/internal/api/common"
	"github.com/SirNacou/weeate/backend/internal/api/foods"
	application "github.com/SirNacou/weeate/backend/internal/app"
	domain "github.com/SirNacou/weeate/backend/internal/domain"
	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/db"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/logger"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/repositories"
)

func main() {
	// Setup logger
	log := logger.NewLogger()

	// Setup context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = log.WithContext(ctx)

	// Setup configuration
	envConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err)
	}

	// Setup Supabase auth
	supabaseClient, err := supabase.NewClient(envConfig.SUPABASE_URL, envConfig.SUPABASE_API_KEY, &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("Failed to initalize the client: ", err)
	}

	// Database connection
	db, err := db.ConnectToPostgres(ctx, envConfig)
	if err != nil {
		log.Fatal().Err(err)
	}

	db.AutoMigrate(&domain.Food{})

	// Setup repositories
	repos := repositories.NewRepositories(db)

	// Setup application handlers
	handlers := application.NewHandlers(&repos, supabaseClient)

	// Setup Fiber app
	app := fiber.New()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log,
	}))

	app.Use(common.CORSMiddleware(envConfig))

	authMiddleware, err := common.AuthMiddleware(ctx, envConfig)
	if err != nil {
		log.Fatal().Err(err)
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

	if err := app.Listen(fmt.Sprintf(":%v", envConfig.PORT)); err != nil {
		log.Fatal().Msgf("Failed to run server: %v", err)
	}
}
