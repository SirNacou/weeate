package foods

import (
	"github.com/SirNacou/weeate/backend/internal/app/foods"
	"github.com/gofiber/fiber/v2"
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

func (e *AddFoodEndpoint) AddFood(c *fiber.Ctx) error {
	result, err := e.addFoodHandler.Handle(c.UserContext(), foods.AddFoodCommand{})
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(AddFoodResponse{
		FoodID: result.FoodID.String(),
	})
}
