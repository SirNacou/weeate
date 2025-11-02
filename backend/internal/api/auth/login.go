package api_auth

import (
	"errors"
	"net/http"

	app_auth "github.com/SirNacou/weeate/backend/internal/app/auth"
	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
	"github.com/SirNacou/weeate/backend/pkg/custom_errors"
	"github.com/labstack/echo/v4"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (e *AuthEndpoint) login(c echo.Context) error {
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, custom_errors.ErrInvalidInput)
	}

	ctx := c.Request().Context()

	token, err := e.loginCH.Handle(ctx, app_auth.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain_auth.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusUnauthorized, domain_auth.ErrInvalidCredentials)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, LoginResponse{Token: token})
}
