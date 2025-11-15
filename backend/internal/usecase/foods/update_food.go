package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type UpdateFoodCommand struct {
	ID          uuid.UUID
	Name        string
	ImageFileId string
	Description string
	Price       int64
}

type UpdateFoodCommandHandler struct {
	db *gorm.DB
}

func NewUpdateFoodCommandHandler(db *gorm.DB) UpdateFoodCommandHandler {
	return UpdateFoodCommandHandler{
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
