package foods

import (
	"context"

	domain "github.com/SirNacou/weeate/backend/internal/domain"
	"github.com/gofrs/uuid/v5"
)

type GetFoodsQuery struct{}

type GetFoodsQueryResult struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
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
	return []GetFoodsQueryResult{}, nil
}
