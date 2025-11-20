package services

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type UpdateFoodCommand struct {
	ID          uuid.UUID `json:"id" format:"uuid" doc:"The ID of the food item to be updated"`
	Name        string    `json:"name" doc:"The name of the food item"`
	ImageFileId string    `json:"image_file_id,omitempty" doc:"The ID of the image file for the food item"`
	Description string    `json:"description" doc:"A description of the food item"`
	Price       int64     `json:"price" doc:"The price of the food item in cents"`
}

type UpdateFoodCommandHandler struct {
	db *gorm.DB
}

func NewUpdateFoodCommandHandler(db *gorm.DB) *UpdateFoodCommandHandler {
	return &UpdateFoodCommandHandler{
		db: db,
	}
}

func (h *UpdateFoodCommandHandler) Handle(ctx context.Context, cmd UpdateFoodCommand) error {
	food := domain.Food{}
	if err := h.db.WithContext(ctx).First(&food, "id = ?", cmd.ID).Error; err != nil {
		return err
	}

	if err := food.UpdateDetails(cmd.Name, cmd.ImageFileId, "", cmd.Description, cmd.Price); err != nil {
		return err
	}

	return h.db.WithContext(ctx).Save(&food).Error
}
