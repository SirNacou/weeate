package services

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/domain"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type GetShoppingOrderQuery struct {
	Date time.Time `json:"date" format:"date" doc:"The date to get the shopping order for"`
}

type GetShoppingOrderResponse struct {
	// Define the response fields here
	Items      []ShoppingItem `json:"items"`
	TotalPrice int64          `json:"total_price"`
}

type ShoppingItem struct {
	FoodID     string             `json:"food_id"`
	FoodName   string             `json:"food_name"`
	Quantity   int64              `json:"quantity"`
	UnitPrice  int64              `json:"unit_price"`
	TotalPrice int64              `json:"total_price"`
	Users      []auth.UserProfile `json:"users"`
}

type GetShoppingOrderQueryHandler struct {
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewGetShoppingOrderQueryHandler(db *gorm.DB, supabaseService *auth.SupabaseService) *GetShoppingOrderQueryHandler {
	return &GetShoppingOrderQueryHandler{
		db:              db,
		supabaseService: supabaseService,
	}
}

func (h *GetShoppingOrderQueryHandler) Handle(ctx context.Context, query GetShoppingOrderQuery) (*GetShoppingOrderResponse, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, auth.ErrUserNotFoundInContext
	}

	order, err := gorm.G[domain.Order](h.db).
		Preload("OrderItems.Details", nil).
		Where("order_date = ? AND buyer_user_id = ?", query.Date, user.ID).First(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &GetShoppingOrderResponse{
		Items: lo.Map(order.OrderItems, func(item domain.OrderItem, _ int) ShoppingItem {
			return ShoppingItem{
				FoodID:     item.FoodID.String(),
				FoodName:   item.FoodName,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				TotalPrice: item.TotalPrice,
				Users:      item.Users, // Adjust if item.Users is not []auth.UserProfile
			}
		}),
		TotalPrice: order.TotalPrice,
	}, nil
}
