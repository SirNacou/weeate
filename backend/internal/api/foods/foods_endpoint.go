package foods

import (
	"github.com/SirNacou/weeate/backend/internal/app"
	"github.com/danielgtaylor/huma/v2"
)

type FoodsEndpoint struct {
	handlers app.Handlers
}

func NewFoodsEndpoint(h app.Handlers) *FoodsEndpoint {
	return &FoodsEndpoint{
		handlers: h,
	}
}

func (e *FoodsEndpoint) Register(api huma.API) error {
	group := huma.NewGroup(api, "/foods")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Foods",
		Description: "Endpoints for managing food items",
	})

	// Register all food endpoints
	NewGetFoodEndpoint(e.handlers.GetFoodsHandler).Register(group)
	NewAddFoodEndpoint(e.handlers.AddFoodHandler).Register(group)
	NewUpdateFoodEndpoint(e.handlers.UpdateFoodHandler).Register(group)
	NewDeleteFoodEndpoint(e.handlers.DeleteFoodHandler).Register(group)

	return nil
}
