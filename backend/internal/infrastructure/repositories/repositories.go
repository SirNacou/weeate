package repositories

import (
	"github.com/SirNacou/weeate/backend/internal/domain"
	"gorm.io/gorm"
)

type Repositories struct {
	// Add repository fields here
	FoodRepo domain.FoodRepository
}

func NewRepositories(db *gorm.DB) Repositories {
	foodRepository := NewGormFoodRepository(db)

	return Repositories{
		FoodRepo: foodRepository,
	}
}
