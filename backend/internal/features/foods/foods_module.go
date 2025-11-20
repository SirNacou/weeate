package foods

import (
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods/services"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterFoodsModule(api huma.API, db *gorm.DB, supabaseService *auth.SupabaseService) {
	getFoodHdl := services.NewGetFoodsQueryHandler(db, supabaseService)
	addFoodHdl := services.NewAddFoodCommandHandler(db)
	updateFoodHdl := services.NewUpdateFoodCommandHandler(db)
	deleteFoodHdl := services.NewDeleteFoodCommandHandler(db)

	foodsEndpoint := NewFoodsEndpoint(
		getFoodHdl,
		addFoodHdl,
		updateFoodHdl,
		deleteFoodHdl,
	)

	group := huma.NewGroup(api, "/foods")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Foods",
		Description: "Endpoints for managing food items",
	})

	huma.Get(group, "/", foodsEndpoint.getFoods)
	huma.Post(group, "/", foodsEndpoint.addFood)
	huma.Put(group, "/{id}", foodsEndpoint.updateFood)
	huma.Delete(group, "/{id}", foodsEndpoint.deleteFood)
}
