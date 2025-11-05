package auth

import (
	"net/http"

	app_auth "github.com/SirNacou/weeate/backend/internal/app/auth"
	"github.com/labstack/echo/v4"
)

type DeleteUserRequest struct {
	userId uint `param:"id"`
}

func (e *AuthEndpoint) deleteUser(c echo.Context) error {
	var req DeleteUserRequest
	if err := c.Bind(&req); err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err := e.deleteUserCH.Handle(c.Request().Context(), app_auth.DeleteUserCommand{UserId: req.userId})
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}
	}

	return c.NoContent(http.StatusOK)
}
