package foods

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/foods/services"
	"github.com/gofrs/uuid/v5"
)

type FoodsEndpoint struct {
	getFoodsHandler   *services.GetFoodsQueryHandler
	addFoodHandler    *services.AddFoodCommandHandler
	updateFoodHandler *services.UpdateFoodCommandHandler
	deleteFoodHandler *services.DeleteFoodCommandHandler
}

func NewFoodsEndpoint(get *services.GetFoodsQueryHandler, add *services.AddFoodCommandHandler, update *services.UpdateFoodCommandHandler, delete *services.DeleteFoodCommandHandler) *FoodsEndpoint {
	return &FoodsEndpoint{
		getFoodsHandler:   get,
		addFoodHandler:    add,
		updateFoodHandler: update,
		deleteFoodHandler: delete,
	}
}

func (e *FoodsEndpoint) getFoods(ctx context.Context, req *struct{ services.GetFoodsQuery }) (*api.Response[[]services.GetFoodsQueryResponse], error) {
	res, err := e.getFoodsHandler.Handle(ctx, req.GetFoodsQuery)
	if err != nil {
		return nil, err
	}

	return api.NewResponse(&res), nil
}

func (e *FoodsEndpoint) addFood(ctx context.Context, req *struct {
	Body   services.AddFoodCommand
},
) (*api.Response[services.AddFoodResult], error) {
	res, err := e.addFoodHandler.Handle(ctx, req.Body)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(res), nil
}

func (e *FoodsEndpoint) updateFood(ctx context.Context, req *struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"The ID of the food item to be updated"`
	Body struct {
		Name        string `json:"name" doc:"The name of the food item"`
		ImageFileId string `json:"image_file_id,omitempty" doc:"The ID of the image file for the food item"`
		Description string `json:"description" doc:"A description of the food item"`
		Price       int64  `json:"price" doc:"The price of the food item in cents"`
	}
},
) (*struct{}, error) {
	err := e.updateFoodHandler.Handle(ctx, services.UpdateFoodCommand{
		ID:          req.ID,
		Name:        req.Body.Name,
		ImageFileId: req.Body.ImageFileId,
		Description: req.Body.Description,
		Price:       req.Body.Price,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (e *FoodsEndpoint) deleteFood(ctx context.Context, req *struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"The ID of the food item to be deleted"`
},
) (*struct{}, error) {
	err := e.deleteFoodHandler.Handle(ctx, services.DeleteFoodCommand{FoodID: req.ID})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
