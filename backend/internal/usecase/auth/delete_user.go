package app_auth

import (
	"context"

	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
)

type DeleteUserCommand struct {
	UserId uint
}

type DeleteUserResponse struct {
	Success bool `json:"success"`
}

type DeleteUserCommandHandler struct {
	userRepo domain_auth.UserRepository
}

func NewDeleteUserCommandHandler(userRepo domain_auth.UserRepository) *DeleteUserCommandHandler {
	return &DeleteUserCommandHandler{
		userRepo: userRepo,
	}
}

func (h *DeleteUserCommandHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	err := h.userRepo.DeleteUser(cmd.UserId)
	if err != nil {
		return err
	}

	return nil
}
