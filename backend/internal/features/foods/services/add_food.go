package services

import (
	"context"
	"fmt"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type AddFoodCommand struct {
	Name        string `json:"name" doc:"The name of the food item"`
	Price       int64  `json:"price" doc:"The price of the food item in cents"`
	Description string `json:"description" doc:"A description of the food item"`
	ImageFileID string `json:"image_file_id,omitempty" doc:"The ID of the image file for the food item"`
}

type AddFoodResult struct {
	FoodID uuid.UUID `json:"food_id" doc:"The ID of the newly created food item"`
}

type AddFoodCommandHandler struct {
	db *gorm.DB
}

func NewAddFoodCommandHandler(db *gorm.DB) *AddFoodCommandHandler {
	return &AddFoodCommandHandler{
		db: db,
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

	if err := h.db.WithContext(ctx).Create(food).Error; err != nil {
		return nil, fmt.Errorf("failed to create food: %w", err)
	}

	return &AddFoodResult{
		food.ID,
	}, nil
}
