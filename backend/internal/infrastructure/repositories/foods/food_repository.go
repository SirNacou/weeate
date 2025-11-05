package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/domain/foods"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type GormFoodRepository struct {
	db *gorm.DB
}

func NewGormFoodRepository(db *gorm.DB) *GormFoodRepository {
	return &GormFoodRepository{db: db}
}

func (r *GormFoodRepository) WithTx(tx *gorm.DB) *GormFoodRepository {
	return &GormFoodRepository{db: tx}
}

func (r *GormFoodRepository) FindByID(ctx context.Context, id uuid.UUID) (foods.Food, error) {
	return gorm.G[foods.Food](r.db).Where(&foods.Food{ID: id}).First(ctx)
}

func (r *GormFoodRepository) FindAllByID(ctx context.Context, ids ...uuid.UUID) ([]foods.Food, error) {
	return gorm.G[foods.Food](r.db).Where("id IN ?", ids).Find(ctx)
}

func (r *GormFoodRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]foods.Food, error) {
	return gorm.G[foods.Food](r.db).Where(&foods.Food{UserID: userID}).Find(ctx)
}

func (r *GormFoodRepository) Create(ctx context.Context, food *foods.Food) error {
	return gorm.G[foods.Food](r.db).Create(ctx, food)
}

func (r *GormFoodRepository) Update(ctx context.Context, food *foods.Food) error {
	_, err := gorm.G[foods.Food](r.db).Where(&foods.Food{ID: food.ID}).Select("*").Updates(ctx, *food)
	return err
}

func (r *GormFoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[foods.Food](r.db).Where(&foods.Food{ID: id}).Delete(ctx)
	return err
}
