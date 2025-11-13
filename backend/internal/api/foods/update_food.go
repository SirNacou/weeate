package foods

import (
	"context"
	"net/http"

	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofrs/uuid/v5"
)

type UpdateFoodRequest struct {
	Name        string `json:"name"`
	ImageFileId string `json:"image_file_id"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type UpdateFoodEndpoint struct {
	updateFoodHandler foods.UpdateFoodCommandHandler
}

func NewUpdateFoodEndpoint(updateFoodHdl foods.UpdateFoodCommandHandler) *UpdateFoodEndpoint {
	return &UpdateFoodEndpoint{
		updateFoodHandler: updateFoodHdl,
	}
}

func (e *UpdateFoodEndpoint) Register(group huma.API) {
	huma.Register(group, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/{id}",
		Summary: "Update food",
	}, e.UpdateFood)
}

func (e *UpdateFoodEndpoint) UpdateFood(ctx context.Context, req *struct {
	ID   uuid.UUID `path:"id"`
	Body UpdateFoodRequest
},
) (*struct{}, error) {
	err := e.updateFoodHandler.Handle(ctx, foods.UpdateFoodCommand{
		ID:          req.ID,
		Name:        req.Body.Name,
		ImageFileId: req.Body.ImageFileId,
		Description: req.Body.Description,
		Price:       req.Body.Price,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("", err)
	}

	return nil, nil
}
