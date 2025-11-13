package foods

import (
	"context"

	domain "github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type GetFoodsQuery struct{}

type GetFoodsQueryResult struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ImageURL        string    `json:"image_url"`
	Description     string    `json:"description"`
	Price           int64     `json:"price"`
	UserID          uuid.UUID `json:"user_id"`
	UserDisplayName string    `json:"username"`
	UserAvatarURL   string    `json:"user_avatar_url"`
}

type GetFoodsQueryHandler struct {
	foodRepo domain.FoodRepository
}

func NewGetFoodsQueryHandler(foodRepo domain.FoodRepository) GetFoodsQueryHandler {
	return GetFoodsQueryHandler{
		foodRepo: foodRepo,
	}
}

func (h *GetFoodsQueryHandler) Handle(ctx context.Context, query GetFoodsQuery) ([]GetFoodsQueryResult, error) {
	res, err := h.foodRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]GetFoodsQueryResult, 0, len(res))
	for _, food := range res {
		results = append(results, GetFoodsQueryResult{
			ID:          food.ID,
			Name:        food.Name,
			ImageURL:    food.ImageURL,
			Description: food.Description,
			Price:       food.Price,
			UserID:      food.UserID,
		})
	}

	return results, nil
}
