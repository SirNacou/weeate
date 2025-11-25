package services

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
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

	foodIDs := lo.Map(order.OrderItems, func(item domain.OrderItem, _ int) string {
		return item.FoodID.String()
	})

	foods, err := gorm.G[food_domain.Food](h.db).
		Where("id IN ?", foodIDs).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	foodMap := lo.KeyBy(foods, func(food food_domain.Food) string {
		return food.ID.String()
	})

	userIDs := lo.FlatMap(order.OrderItems, func(item domain.OrderItem, _ int) []string {
		return lo.Map(item.Details, func(detail domain.OrderItemDetail, _ int) string {
			return detail.UserID
		})
	})

	userProfiles, err := h.supabaseService.GetUserProfilesByIDs(userIDs...)
	if err != nil {
		return nil, err
	}
	userProfileMap := lo.KeyBy(userProfiles, func(profile auth.UserProfile) string {
		return profile.ID.String()
	})

	return &GetShoppingOrderResponse{
		Items: lo.Map(order.OrderItems, func(item domain.OrderItem, _ int) ShoppingItem {
			quantity := lo.SumBy(item.Details, func(detail domain.OrderItemDetail) int64 {
				return detail.Quantity
			})
			return ShoppingItem{
				FoodID:     item.FoodID.String(),
				FoodName:   foodMap[item.FoodID.String()].Name,
				Quantity:   quantity,
				UnitPrice:  item.PriceAtOrder,
				TotalPrice: item.PriceAtOrder * quantity,
				Users: lo.Map(item.Details, func(detail domain.OrderItemDetail, _ int) auth.UserProfile {
					return userProfileMap[detail.UserID]
				}),
			}
		}),
		TotalPrice: order.TotalPrice,
	}, nil
}
