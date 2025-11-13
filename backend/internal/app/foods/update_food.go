package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type UpdateFoodCommand struct {
	ID          uuid.UUID
	Name        string
	ImageFileId string
	Description string
	Price       int64
}

type UpdateFoodCommandHandler struct {
	foodRepo domain.FoodRepository
}

func NewUpdateFoodCommandHandler(foodRepo domain.FoodRepository) UpdateFoodCommandHandler {
	return UpdateFoodCommandHandler{
		foodRepo: foodRepo,
	}
}

func (h *UpdateFoodCommandHandler) Handle(ctx context.Context, cmd UpdateFoodCommand) error {
	food, err := h.foodRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if err = food.UpdateDetails(cmd.Name, cmd.ImageFileId, "", cmd.Description, cmd.Price); err != nil {
		return err
	}

	return h.foodRepo.Update(ctx, &food)
}
