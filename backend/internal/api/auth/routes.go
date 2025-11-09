package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthEndpoint(f *fiber.App, endpoint AuthEndpoint) error {
	auth := f.Group("/auth")
	auth.Delete("users/:id", endpoint.DeleteUser)
	return nil
}
