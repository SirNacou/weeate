package repositories

import (
	"context"

	domain "github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type GormFoodRepository struct {
	db *gorm.DB
}

func NewGormFoodRepository(db *gorm.DB) domain.FoodRepository {
	return &GormFoodRepository{db: db}
}

func (r *GormFoodRepository) WithTx(tx *gorm.DB) domain.FoodRepository {
	return &GormFoodRepository{db: tx}
}

func (r *GormFoodRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Food, error) {
	return gorm.G[domain.Food](r.db).Where(&domain.Food{ID: id}).First(ctx)
}

func (r *GormFoodRepository) FindAllByID(ctx context.Context, ids ...uuid.UUID) ([]domain.Food, error) {
	return gorm.G[domain.Food](r.db).Where("id IN ?", ids).Find(ctx)
}

func (r *GormFoodRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Food, error) {
	return gorm.G[domain.Food](r.db).Where(&domain.Food{UserID: userID}).Find(ctx)
}

func (r *GormFoodRepository) Create(ctx context.Context, food *domain.Food) error {
	return gorm.G[domain.Food](r.db).Create(ctx, food)
}

func (r *GormFoodRepository) Update(ctx context.Context, food *domain.Food) error {
	_, err := gorm.G[domain.Food](r.db).Where(&domain.Food{ID: food.ID}).Select("*").Updates(ctx, *food)
	return err
}

func (r *GormFoodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := gorm.G[domain.Food](r.db).Where(&domain.Food{ID: id}).Delete(ctx)
	return err
}
