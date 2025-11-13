package common

import (
	"net/http"
	"strings"

	config "github.com/SirNacou/weeate/backend/internal/infrastructure/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware(cfg config.Config) fiber.Handler {
	if cfg.GO_ENV == config.EnvDevelopment {
		return cors.New(cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return true // Allow all origins in development
			},
			AllowMethods:     strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead}, ", "),
			AllowHeaders:     strings.Join([]string{fiber.HeaderOrigin, fiber.HeaderContentType, fiber.HeaderAccept, fiber.HeaderAuthorization}, ", "),
			AllowCredentials: true,
		})
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join([]string{"https://weeate.nacou.uk"}, ", "),
		AllowMethods:     strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions, http.MethodHead}, ", "),
		AllowHeaders:     strings.Join([]string{fiber.HeaderOrigin, fiber.HeaderContentType, fiber.HeaderAccept, fiber.HeaderAuthorization}, ", "),
		AllowCredentials: true,
	})
}
