package usecase

import (
	"github.com/SirNacou/weeate/backend/internal/usecase/foods"
	"github.com/supabase-community/supabase-go"
	"gorm.io/gorm"
)

type Handlers struct {
	// Add common fields for handlers if needed
	GetFoodsHandler   foods.GetFoodsQueryHandler
	AddFoodHandler    foods.AddFoodCommandHandler
	UpdateFoodHandler foods.UpdateFoodCommandHandler
	DeleteFoodHandler foods.DeleteFoodCommandHandler
}

func NewHandlers(db *gorm.DB, supabaseClient *supabase.Client) Handlers {
	getFoodHdl := foods.NewGetFoodsQueryHandler(db, supabaseClient)
	addFoodHdl := foods.NewAddFoodCommandHandler(db)
	updateFoodHdl := foods.NewUpdateFoodCommandHandler(db)
	deleteFoodHdl := foods.NewDeleteFoodCommandHandler(db)
	return Handlers{
		GetFoodsHandler:   getFoodHdl,
		AddFoodHandler:    addFoodHdl,
		UpdateFoodHandler: updateFoodHdl,
		DeleteFoodHandler: deleteFoodHdl,
	}
}
