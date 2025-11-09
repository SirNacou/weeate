package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	// Database connection
	db, err := db.ConnectToPostgres(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&domain.Food{})

	// Initialize Fiber app
	app := fiber.New()

	app.Use(logger.New(logger.ConfigDefault))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join([]string{"http://localhost:3000", "http://localhost:3001", "https://weeate.nacou.uk"}, ", "),
		AllowMethods:     strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}, ", "),
		AllowHeaders:     strings.Join([]string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}, ", "),
		AllowCredentials: true,
	}))

	authware, err := auth.NewAuthMiddleware(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}
	app.Use(authware.Handle)

	repos := repositories.NewRepositories(db)
	handlers := application.NewHandlers(&repos)
	foodsEndpoint := foods.NewFoodsEndpoint(handlers)

	foodsEndpoint.Register(app)

	app.Get("/", func(ctx *fiber.Ctx) error {
		user, err := auth.GetUserContext(ctx)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).SendString(err.Error())
		}
		ctx.Status(http.StatusOK).JSON(user)
		return ctx.Send(ctx.Body())
	})

	if err := app.Listen(fmt.Sprintf(":%v", config.PORT)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
