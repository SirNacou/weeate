package auth

import (
	"net/http"

	app_auth "github.com/SirNacou/weeate/backend/internal/usecase/auth"
	"github.com/gofiber/fiber/v2"
)

type DeleteUserRequest struct {
	userId string `uri:"id" binding:"required"`
}

func (e *AuthEndpoint) DeleteUser(c *fiber.Ctx) error {
	if err := e.deleteUserCH.Handle(c.Context(), app_auth.DeleteUserCommand{UserId: uint(0)}); err != nil {
		return c.Status(http.StatusBadRequest).Send([]byte(err.Error()))
	}

	return c.SendStatus(http.StatusOK)
}
