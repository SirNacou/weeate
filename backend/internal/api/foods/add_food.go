package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/danielgtaylor/huma/v2"
)

type AddFoodRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddFoodResponse struct {
	FoodID string `json:"food_id"`
}

type AddFoodEndpoint struct {
	addFoodHandler foods.AddFoodCommandHandler
}

func NewAddFoodEndpoint(addFoodHdl foods.AddFoodCommandHandler) *AddFoodEndpoint {
	return &AddFoodEndpoint{
		addFoodHandler: addFoodHdl,
	}
}

func (e *AddFoodEndpoint) AddFood(ctx context.Context, req *struct {
	AddFoodRequest
},
) (*AddFoodResponse, error) {
	result, err := e.addFoodHandler.Handle(ctx, foods.AddFoodCommand{})
	if err != nil {
		return nil, huma.Error400BadRequest("", err)
	}

	return &AddFoodResponse{
		FoodID: result.FoodID.String(),
	}, nil
}
