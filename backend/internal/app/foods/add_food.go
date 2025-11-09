package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type AddFoodCommand struct {
	UserID string
}

type AddFoodResult struct {
	FoodID uuid.UUID
}

type AddFoodCommandHandler struct {
	foodRepo domain.FoodRepository
}

func NewAddFoodCommandHandler(foodRepo domain.FoodRepository) AddFoodCommandHandler {
	return AddFoodCommandHandler{
		foodRepo: foodRepo,
	}
}

func (h *AddFoodCommandHandler) Handle(ctx context.Context, command AddFoodCommand) (AddFoodResult, error) {
	return AddFoodResult{}, nil
}
