package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/api"
	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
)

type GetFoodsRequest struct{}

type GetFoodsResponse struct {
	Result []GetFoodsResponseItem `json:"result"`
}

type GetFoodsResponseItem struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
}

type GetFoodEndpoint struct {
	getFoodsHandler foods.GetFoodsQueryHandler
}

func NewGetFoodEndpoint(getFoodsHdl foods.GetFoodsQueryHandler) *GetFoodEndpoint {
	return &GetFoodEndpoint{
		getFoodsHandler: getFoodsHdl,
	}
}

func (e *GetFoodEndpoint) GetFoods(ctx context.Context, req *struct{}) (*api.Response[GetFoodsResponse], error) {
	r, err := e.getFoodsHandler.Handle(ctx, foods.GetFoodsQuery{})
	if err != nil {
		return nil, huma.Error400BadRequest("", err)
	}

	var res []GetFoodsResponseItem
	for _, food := range r {
		res = append(res, GetFoodsResponseItem{
			ID:          food.ID,
			Name:        food.Name,
			ImageURL:    food.ImageURL,
			Description: food.Description,
			Price:       food.Price,
			UserID:      food.UserID,
		})
	}

	return api.NewResponse(GetFoodsResponse{
		res,
	}), nil
}
