package foods

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type FoodRepository interface {
	WithTx(tx *gorm.DB) FoodRepository
	FindByID(ctx context.Context, id uuid.UUID) (Food, error)
	FindAllByID(ctx context.Context, ids ...uuid.UUID) ([]Food, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]Food, error)
	Create(ctx context.Context, food *Food) error
	Update(ctx context.Context, food *Food) error
	Delete(ctx context.Context, id uuid.UUID) error
}
