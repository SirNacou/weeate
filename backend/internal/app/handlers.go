package app

import (
	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/SirNacou/weeate/backend/internal/infrastructure/repositories"
	"github.com/supabase-community/supabase-go"
)

type Handlers struct {
	// Add common fields for handlers if needed
	GetFoodsHandler   foods.GetFoodsQueryHandler
	AddFoodHandler    foods.AddFoodCommandHandler
	UpdateFoodHandler foods.UpdateFoodCommandHandler
	DeleteFoodHandler foods.DeleteFoodCommandHandler
}

func NewHandlers(repos *repositories.Repositories, supabaseClient *supabase.Client) Handlers {
	getFoodHdl := foods.NewGetFoodsQueryHandler(repos.FoodRepo, supabaseClient)
	addFoodHdl := foods.NewAddFoodCommandHandler(repos.FoodRepo)
	updateFoodHdl := foods.NewUpdateFoodCommandHandler(repos.FoodRepo)
	deleteFoodHdl := foods.NewDeleteFoodCommandHandler(repos.FoodRepo)
	return Handlers{
		GetFoodsHandler:   getFoodHdl,
		AddFoodHandler:    addFoodHdl,
		UpdateFoodHandler: updateFoodHdl,
		DeleteFoodHandler: deleteFoodHdl,
	}
}
