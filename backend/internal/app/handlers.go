package app

import (
	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/repositories"
)

type Handlers struct {
	// Add common fields for handlers if needed
	AddFoodHandler    foods.AddFoodCommandHandler
	UpdateFoodHandler foods.UpdateFoodCommandHandler
	GetFoodsHandler   foods.GetFoodsQueryHandler
}

func NewHandlers(repos *repositories.Repositories) Handlers {
	addFoodHdl := foods.NewAddFoodCommandHandler(repos.FoodRepo)
	updateFoodHdl := foods.NewUpdateFoodCommandHandler(repos.FoodRepo)
	getFoodHdl := foods.NewGetFoodsQueryHandler(repos.FoodRepo)
	return Handlers{
		AddFoodHandler:    addFoodHdl,
		UpdateFoodHandler: updateFoodHdl,
		GetFoodsHandler:   getFoodHdl,
	}
}
