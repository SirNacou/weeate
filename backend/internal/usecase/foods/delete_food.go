package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type DeleteFoodCommand struct {
	FoodID uuid.UUID
}

type DeleteFoodCommandHandler struct {
	db *gorm.DB
}

func NewDeleteFoodCommandHandler(db *gorm.DB) DeleteFoodCommandHandler {
	return DeleteFoodCommandHandler{
		db: db,
	}
}

func (h *DeleteFoodCommandHandler) Handle(ctx context.Context, command DeleteFoodCommand) error {
	if err := h.db.WithContext(ctx).Delete(&domain.Food{}, "id = ?", command.FoodID).Error; err != nil {
		return err
	}

	return nil
}
