package foods

import (
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods/add_food"
	"github.com/SirNacou/weeate/backend/internal/features/foods/delete_food"
	"github.com/SirNacou/weeate/backend/internal/features/foods/get_foods"
	"github.com/SirNacou/weeate/backend/internal/features/foods/update_food"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterFoodsModule(api huma.API, db *gorm.DB, supabaseService *auth.SupabaseService) {
	getFoodHdl := get_foods.NewGetFoodsQueryHandler(db, supabaseService)
	addFoodHdl := add_food.NewAddFoodCommandHandler(db)
	updateFoodHdl := update_food.NewUpdateFoodCommandHandler(db)
	deleteFoodHdl := delete_food.NewDeleteFoodCommandHandler(db)

	group := huma.NewGroup(api, "/foods")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Foods",
		Description: "Endpoints for managing food items",
	})

	get_foods.NewGetFoodEndpoint(getFoodHdl).Register(group)
	add_food.NewAddFoodEndpoint(addFoodHdl).Register(group)
	update_food.NewUpdateFoodEndpoint(updateFoodHdl).Register(group)
	delete_food.NewDeleteFoodEndpoint(deleteFoodHdl).Register(group)
}
