package foods

import (
	"context"

	domain "github.com/SirNacou/weeate/backend/internal/domain"
	user_domain "github.com/SirNacou/weeate/backend/internal/domain/auth"
	"github.com/gofrs/uuid/v5"
)

type GetFoodsByUserIDQuery struct {
	UserID string
}

type GetFoodsByUserIDQueryResult struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
}

type GetFoodsByUserIDQueryHandler struct {
	foodRepo domain.FoodRepository
}

func NewGetFoodsByUserIDQueryHandler(foodRepo domain.FoodRepository) *GetFoodsByUserIDQueryHandler {
	return &GetFoodsByUserIDQueryHandler{
		foodRepo: foodRepo,
	}
}

func (h *GetFoodsByUserIDQueryHandler) Handle(ctx context.Context, query GetFoodsByUserIDQuery) ([]GetFoodsByUserIDQueryResult, error) {
	userID, err := uuid.FromString(query.UserID)
	if err != nil {
		return nil, user_domain.ErrUserNotFound
	}

	foodList, err := h.foodRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]GetFoodsByUserIDQueryResult, len(foodList))
	for _, v := range foodList {
		result = append(result, GetFoodsByUserIDQueryResult{
			ID:          v.ID,
			Name:        v.Name,
			ImageURL:    v.ImageURL,
			Description: v.Description,
			Price:       v.Price,
			UserID:      v.UserID,
		})
	}

	return result, err
}
