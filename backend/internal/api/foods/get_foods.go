package foods

import (
	"context"
	"net/http"

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
	ID          uuid.UUID         `json:"id"`
	Name        string            `json:"name"`
	ImageURL    string            `json:"image_url"`
	Description string            `json:"description"`
	Price       int64             `json:"price"`
	User        foods.UserProfile `json:"user"`
}

type GetFoodEndpoint struct {
	getFoodsHandler foods.GetFoodsQueryHandler
}

func NewGetFoodEndpoint(getFoodsHdl foods.GetFoodsQueryHandler) *GetFoodEndpoint {
	return &GetFoodEndpoint{
		getFoodsHandler: getFoodsHdl,
	}
}

func (e *GetFoodEndpoint) Register(group huma.API) {
	huma.Register(group, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/",
		Summary:       "Get List Foods",
		Description:   "Retrieve a list of foods.",
		DefaultStatus: http.StatusOK,
	}, e.GetFoods)
}

func (e *GetFoodEndpoint) GetFoods(ctx context.Context, req *struct{}) (*api.Response[GetFoodsResponse], error) {
	r, err := e.getFoodsHandler.Handle(ctx, foods.GetFoodsQuery{})
	if err != nil {
		return nil, huma.Error400BadRequest("", err)
	}

	res := make([]GetFoodsResponseItem, 0, len(r))
	for _, food := range r {
		res = append(res, GetFoodsResponseItem{
			ID:          food.ID,
			Name:        food.Name,
			ImageURL:    food.ImageURL,
			Description: food.Description,
			Price:       food.Price,
			User:        food.User,
		})
	}

	return api.NewResponse(GetFoodsResponse{
		res,
	}), nil
}
