package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type DeleteFoodCommand struct {
	FoodID uuid.UUID
}

type DeleteFoodCommandHandler struct {
	// Add necessary dependencies here, e.g., food repository
	foodRepo domain.FoodRepository
}

func NewDeleteFoodCommandHandler(foodRepo domain.FoodRepository) DeleteFoodCommandHandler {
	return DeleteFoodCommandHandler{
		foodRepo: foodRepo,
	}
}

func (h *DeleteFoodCommandHandler) Handle(ctx context.Context, command DeleteFoodCommand) error {
	if err := h.foodRepo.Delete(ctx, command.FoodID); err != nil {
		return err
	}

	return nil
}
