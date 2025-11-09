package foods

import (
	"github.com/SirNacou/weeate/backend/internal/app"
	"github.com/gofiber/fiber/v2"
)

type FoodsEndpoint struct {
	handlers app.Handlers
}

func NewFoodsEndpoint(h app.Handlers) *FoodsEndpoint {
	return &FoodsEndpoint{
		handlers: h,
	}
}

func (e *FoodsEndpoint) Register(f *fiber.App) error {
	group := f.Group("/foods")

	getFoodsEndpoint := NewGetFoodEndpoint(e.handlers.GetFoodsHandler)
	group.Get("/", getFoodsEndpoint.GetFoods)

	addFoodEndpoint := NewAddFoodEndpoint(e.handlers.AddFoodHandler)
	group.Post("/", addFoodEndpoint.AddFood)

	return nil
}
