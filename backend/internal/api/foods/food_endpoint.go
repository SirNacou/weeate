package foods

import (
	"net/http"

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

	getFoodsEndpoint := NewGetFoodEndpoint(e.handlers.GetFoodsHandler)
	huma.Register(group, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/",
		Summary:       "Get List Foods",
		Description:   "Retrieve a list of foods.",
		DefaultStatus: http.StatusOK,
	}, getFoodsEndpoint.GetFoods)

	addFoodEndpoint := NewAddFoodEndpoint(e.handlers.AddFoodHandler)
	huma.Register(group, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/",
		Summary:       "Add new food",
		DefaultStatus: http.StatusOK,
	}, addFoodEndpoint.AddFood)

	return nil
}
