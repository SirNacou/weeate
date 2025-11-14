package foods

import (
	"context"
	"net/http"

	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
)

type DeleteFoodEndpoint struct {
	deleteFoodHandler foods.DeleteFoodCommandHandler
}

func NewDeleteFoodEndpoint(deleteFoodHdl foods.DeleteFoodCommandHandler) *DeleteFoodEndpoint {
	return &DeleteFoodEndpoint{
		deleteFoodHandler: deleteFoodHdl,
	}
}

func (e *DeleteFoodEndpoint) Register(group huma.API) {
	huma.Register(group, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete a food item by its ID",
	}, e.DeleteFood)
}

func (e *DeleteFoodEndpoint) DeleteFood(ctx context.Context, req *struct {
	ID uuid.UUID `path:"id" format:"uuid"`
},
) (*struct{}, error) {
	err := e.deleteFoodHandler.Handle(ctx, foods.DeleteFoodCommand{
		FoodID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
