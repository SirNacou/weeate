package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"

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

	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "https://weeate.nacou.uk"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: func(t *jwt.Token) (any, error) {
			ctx := context.Background()
			iss, err := t.Claims.GetIssuer()
			if err != nil {
				return nil, err
			}

			pubKeyUrl, err := url.JoinPath(iss, ".well-known/jwks.json")
			if err != nil {
				return nil, err
			}

			keyID, ok := t.Header["kid"].(string)
			if !ok {
				return nil, errors.New("expecting JWT header to have a key ID in the kid field")
			}

			set, err := jwk.Fetch(ctx, pubKeyUrl)
			if err != nil {
				return nil, err
			}

			key, found := set.LookupKeyID(keyID)
			if !found {
				return nil, fmt.Errorf("unable to find key %q", keyID)
			}

			k, err := key.PublicKey()
			if err != nil {
				return nil, errors.New("unable to get public key")
			}

			var ecdsaK ecdsa.PublicKey
			err = k.Raw(&ecdsaK)
			return &ecdsaK, err
		},
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
