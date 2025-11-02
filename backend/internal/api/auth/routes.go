package api_auth

import (
	"github.com/labstack/echo/v4"
)

func RegisterAuthEndpoint(e *echo.Echo, endpoint AuthEndpoint) error {
	auth := e.Group("/auth")
	auth.POST("/register", endpoint.register)
	auth.POST("/login", endpoint.login)
	auth.DELETE("users/:id", endpoint.deleteUser)
	return nil
}
