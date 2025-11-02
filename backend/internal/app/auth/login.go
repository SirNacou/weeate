package app_auth

import (
	"context"

	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

type LoginCommand struct {
	Username string
	Password string
}

type LoginCommandHandler struct {
	userRepo domain_auth.UserRepository
}

func NewLoginCommandHandler(userRepo domain_auth.UserRepository) *LoginCommandHandler {
	return &LoginCommandHandler{
		userRepo: userRepo,
	}
}

func (h *LoginCommandHandler) Handle(ctx context.Context, cmd LoginCommand) (string, error) {
	user, err := h.userRepo.GetUserByUsername(cmd.Username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cmd.Password)); err != nil {
		return "", domain_auth.ErrInvalidCredentials
	}
	return "", nil
}
