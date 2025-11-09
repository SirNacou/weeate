package foods

import (
	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
)

type GetFoodsRequest struct{}

type GetFoodsResponse struct {
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

func (e *GetFoodEndpoint) GetFoods(ctx *fiber.Ctx) error {
	r, err := e.getFoodsHandler.Handle(ctx.UserContext(), foods.GetFoodsQuery{})
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(r)
}
