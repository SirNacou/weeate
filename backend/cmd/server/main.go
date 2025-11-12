package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"

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

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := flag.Bool("debug", false, "Enable debug mode with more verbose logging")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i any) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i any) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i any) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	// Setup configuration
	envConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err)
	}

	// Database connection
	db, err := db.ConnectToPostgres(ctx, envConfig)
	if err != nil {
		log.Fatal().Err(err)
	}

	db.AutoMigrate(&domain.Food{})

	app := fiber.New()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log,
	}))

	if envConfig.GO_ENV == config.EnvDevelopment {
		log.Println("Running in development mode")
		app.Use(cors.New(cors.ConfigDefault, cors.Config{
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
		log.Fatal().Err(err)
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
		log.Fatal().Msgf("Failed to run server: %v", err)
	}
}
