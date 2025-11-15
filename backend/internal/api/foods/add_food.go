package foods

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/SirNacou/weeate/backend/internal/api"
	"github.com/SirNacou/weeate/backend/internal/usecase/foods"
	"github.com/danielgtaylor/huma/v2"
)

// StringInt is a custom type that can unmarshal from both string and int
type StringInt int

func (si *StringInt) UnmarshalText(text []byte) error {
	val, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	*si = StringInt(val)
	return nil
}

type AddFoodRequest struct {
	Name        string `json:"name" minLength:"1"`
	Price       int64  `json:"price" minimum:"0" multipleOf:"1000"`
	Description string `json:"description" required:"false"`
	ImageFileID string `json:"image_file_id" required:"false"`
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

func (e *AddFoodEndpoint) Register(group huma.API) {
	huma.Register(group, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/",
		Summary:       "Add new food",
		DefaultStatus: http.StatusOK,
	}, e.AddFood)
}

func (e *AddFoodEndpoint) AddFood(ctx context.Context, req *struct{ Body AddFoodRequest }) (*api.Response[AddFoodResponse], error) {
	log.Printf("Received AddFood request: %+v", req)

	result, err := e.addFoodHandler.Handle(ctx, foods.AddFoodCommand{
		Name:        req.Body.Name,
		Price:       req.Body.Price,
		Description: req.Body.Description,
		ImageFileID: req.Body.ImageFileID,
	})
	if err != nil {
		return nil, huma.Error400BadRequest("", err)
	}

	return api.NewResponse(AddFoodResponse{
		FoodID: result.FoodID.String(),
	}), nil
}
