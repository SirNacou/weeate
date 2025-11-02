package main

import (
	"context"
	"log"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v4"

	api_auth "github.com/SirNacou/weeate/backend/internal/api/auth"
	app_auth "github.com/SirNacou/weeate/backend/internal/app/auth"
	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/db"
	infra_auth "github.com/SirNacou/weeate/backend/internal/infrastructure/repositories/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	e := echo.New()

	db, err := db.ConnectToPostgres(ctx, cfg)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&domain_auth.User{})

	userRepo := infra_auth.NewUserRepository(db)

	registerCH := app_auth.NewRegisterCommandHandler(userRepo)
	loginCH := app_auth.NewLoginCommandHandler(userRepo)
	deleteUserCH := app_auth.NewDeleteUserCommandHandler(userRepo)

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "https://weeate.nacou.uk"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: "05d6bf6b-7b42-4d50-9af7-8e9ce634d1c2",
		SigningMethod: "ES256",
		TokenLookup: "header:Authorization:Bearer ",
	}))

	authEndpoint := api_auth.NewAuthEndpoint(*registerCH, *loginCH, *deleteUserCH)

	err = api_auth.RegisterAuthEndpoint(e, *authEndpoint)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("user"))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
