package app_auth

import (
	"context"

	domain_auth "github.com/SirNacou/weeate/backend/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

type RegisterCommand struct {
	Username string
	Password string
	Fullname string
}

type RegisterCommandHandler struct {
	userRepo domain_auth.UserRepository
}

func NewRegisterCommandHandler(userRepo domain_auth.UserRepository) *RegisterCommandHandler {
	return &RegisterCommandHandler{
		userRepo: userRepo,
	}
}

func (h *RegisterCommandHandler) Handle(ctx context.Context, cmd RegisterCommand) error {
	// Implement registration logic here
	if len(cmd.Username) < 3 || len(cmd.Username) > 30 {
		return domain_auth.ErrInvalidUsernameLength
	}
	if len(cmd.Password) < 8 {
		return domain_auth.ErrInvalidPasswordLength
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &domain_auth.User{
		Username: cmd.Username,
		Password: string(hashedPassword),
		FullName: cmd.Fullname,
	}

	return h.userRepo.CreateUser(newUser)
}
