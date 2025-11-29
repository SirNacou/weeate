package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/foods/services"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type FoodsModule struct {
	endpoint *FoodsEndpoint
}

func NewFoodsModule(db *gorm.DB, supabaseService *auth.SupabaseService) *FoodsModule {
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

	return &FoodsModule{
		endpoint: foodsEndpoint,
	}
}

func (m *FoodsModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/foods")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Foods",
		Description: "Endpoints for managing food items",
	})

	huma.Get(group, "/", m.endpoint.getFoods)
	huma.Post(group, "/", m.endpoint.addFood)
	huma.Put(group, "/{id}", m.endpoint.updateFood)
	huma.Delete(group, "/{id}", m.endpoint.deleteFood)
}

func (m *FoodsModule) Jobs() []func(context.Context) {
	return nil
}
