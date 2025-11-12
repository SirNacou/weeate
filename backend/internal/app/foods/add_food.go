package foods

import (
	"context"
	"fmt"

	"github.com/SirNacou/weeate/backend/internal/api/auth"
	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type AddFoodCommand struct {
	Name        string
	Price       int64
	Description string
	ImageFileID string
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

func (h *AddFoodCommandHandler) Handle(ctx context.Context, command AddFoodCommand) (*AddFoodResult, error) {
	user, err := auth.GetUserContext(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.FromString(user.ID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate image id and get ImageURL

	food, err := domain.NewFood(command.Name, "", "", command.Description, command.Price, userID)
	if err != nil {
		return nil, err
	}

	if err := h.foodRepo.Create(ctx, food); err != nil {
		return nil, fmt.Errorf("failed to create food: %w", err)
	}

	return &AddFoodResult{
		food.ID,
	}, nil
}
