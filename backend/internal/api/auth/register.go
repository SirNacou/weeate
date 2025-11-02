package api_auth

import (
	"net/http"

	app_auth "github.com/SirNacou/weeate/backend/internal/app/auth"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8"`
	Fullname string `json:"fullname" validate:"max=100"`
}

func (e *AuthEndpoint) register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	err := e.registerCH.Handle(c.Request().Context(), app_auth.RegisterCommand{
		Username: req.Username,
		Password: req.Password,
		Fullname: req.Fullname,
	})
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: map[string]string{"error": "Failed to register user"},
		}
	}

	return c.NoContent(http.StatusCreated)
}
