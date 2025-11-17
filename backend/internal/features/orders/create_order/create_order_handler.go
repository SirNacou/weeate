package create_order

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
)

type CreateOrderCommand struct {
	// Define request fields here
}

type CreateOrderCommandHandler struct{}

func NewCreateOrderCommandHandler() *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{}
}

func (h *CreateOrderCommandHandler) Handle(ctx context.Context, req *CreateOrderCommand) error {
	// Implementation goes here
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return auth.ErrUserNotFoundInContext
	}
	

	return nil
}
