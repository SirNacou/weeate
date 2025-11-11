package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/SirNacou/weeate/backend/internal/api/auth"
	"github.com/SirNacou/weeate/backend/internal/api/foods"
	application "github.com/SirNacou/weeate/backend/internal/app"
	domain "github.com/SirNacou/weeate/backend/internal/domain"
	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/db"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/repositories"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup configuration
	envConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Database connection
	db, err := db.ConnectToPostgres(ctx, envConfig)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&domain.Food{})

	app := fiber.New()

	app.Use(logger.New(logger.ConfigDefault))

	if envConfig.GO_ENV == config.EnvDevelopment {
		log.Println("Running in development mode")
		app.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Join([]string{"http://localhost:3000", "http://127.0.0.1:3000"}, ", "),
			AllowMethods:     strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}, ", "),
			AllowHeaders:     strings.Join([]string{fiber.HeaderOrigin, fiber.HeaderContentType, fiber.HeaderAccept, echo.HeaderAuthorization}, ", "),
			AllowCredentials: true,
		}))
	} else {
		log.Println("Running in production mode")
		app.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Join([]string{"https://weeate.nacou.uk"}, ", "),
			AllowMethods:     strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}, ", "),
			AllowHeaders:     strings.Join([]string{fiber.HeaderOrigin, fiber.HeaderContentType, fiber.HeaderAccept, echo.HeaderAuthorization}, ", "),
			AllowCredentials: true,
		}))
	}

	authware, err := auth.NewAuthMiddleware(ctx, envConfig)
	if err != nil {
		log.Fatalln(err)
	}
	app.Use(authware.Handle)

	api := humafiber.New(app, huma.DefaultConfig("Weeate API", "v1.0.0"))

	repos := repositories.NewRepositories(db)
	handlers := application.NewHandlers(&repos)
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
		log.Fatalf("Failed to run server: %v", err)
	}
}
