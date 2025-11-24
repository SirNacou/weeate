package services

import (
	"context"
	"log/slog"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	domain_food "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/orders/domain"
	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

type (
	GetTodayOrdersResponse struct {
		PollID     uuid.UUID        `json:"poll_id"`
		Buyer      auth.UserProfile `json:"buyer"`
		TotalPrice int64            `json:"total_price"`
		OrderItems []OrderItem      `json:"order_items"`
	}
	OrderItem struct {
		FoodName     string            `json:"food_name"`
		FoodImageUrl *string           `json:"food_image_url"`
		PriceAtOrder int64             `json:"price_at_order"`
		Details      []OrderItemDetail `json:"details"`
	}

	OrderItemDetail struct {
		User     auth.UserProfile `json:"user"`
		Quantity int64            `json:"quantity"`
	}

	GetTodayOrdersQueryHandler struct {
		db              *gorm.DB
		supabaseService *auth.SupabaseService
	}
)

func NewGetTodayOrdersQueryHandler(db *gorm.DB, supabaseService *auth.SupabaseService) *GetTodayOrdersQueryHandler {
	return &GetTodayOrdersQueryHandler{
		db:              db,
		supabaseService: supabaseService,
	}
}

func (h *GetTodayOrdersQueryHandler) Handle(ctx context.Context, query *struct{}) ([]GetTodayOrdersResponse, error) {
	orders, err := gorm.G[domain.Order](h.db).
		Preload("OrderItems.Details", nil).
		Where("order_date = CURRENT_DATE").
		Find(ctx)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return []GetTodayOrdersResponse{}, nil
	}

	userIDs, foodIDs := collectIDs(orders)

	userProfileMap, foodMap, err := h.fetchData(ctx, userIDs, foodIDs)
	if err != nil {
		return nil, err
	}

	return buildResponse(ctx, orders, userProfileMap, foodMap), nil
}

func collectIDs(orders []domain.Order) (userIDs []string, foodIDs []uuid.UUID) {
	userIDMap := make(map[string]struct{})
	foodIDMap := make(map[uuid.UUID]struct{})

	for _, order := range orders {
		userIDMap[order.BuyerUserID] = struct{}{}
		for _, item := range order.OrderItems {
			foodIDMap[item.FoodID] = struct{}{}
			for _, detail := range item.Details {
				userIDMap[detail.UserID] = struct{}{}
			}
		}
	}

	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	for id := range foodIDMap {
		foodIDs = append(foodIDs, id)
	}

	return userIDs, foodIDs
}

func (h *GetTodayOrdersQueryHandler) fetchData(ctx context.Context, userIDs []string, foodIDs []uuid.UUID) (map[string]auth.UserProfile, map[uuid.UUID]domain_food.Food, error) {
	users, err := h.supabaseService.GetUserProfilesByIDs(userIDs)
	if err != nil {
		return nil, nil, err
	}
	userProfileMap := make(map[string]auth.UserProfile, len(users))
	for _, user := range users {
		userProfileMap[user.ID.String()] = user
	}

	foods, err := gorm.G[domain_food.Food](h.db).Where("id IN ?", foodIDs).Find(ctx)
	if err != nil {
		return nil, nil, err
	}
	foodMap := make(map[uuid.UUID]domain_food.Food, len(foods))
	for _, food := range foods {
		foodMap[food.ID] = food
	}

	return userProfileMap, foodMap, nil
}

func buildResponse(ctx context.Context, orders []domain.Order, userProfileMap map[string]auth.UserProfile, foodMap map[uuid.UUID]domain_food.Food) []GetTodayOrdersResponse {
	var response []GetTodayOrdersResponse
	for _, order := range orders {
		buyerProfile, ok := userProfileMap[order.BuyerUserID]
		if !ok {
			slog.ErrorContext(ctx, "buyer not found",
				"buyerID", order.BuyerUserID,
				"orderID", order.ID.String())
			continue
		}

		var orderItems []OrderItem
		for _, item := range order.OrderItems {
			food, ok := foodMap[item.FoodID]
			if !ok {
				slog.ErrorContext(ctx, "food not found",
					"foodID", item.FoodID,
					"orderID", order.ID.String())
				continue
			}

			var details []OrderItemDetail
			for _, detail := range item.Details {
				userProfile, ok := userProfileMap[detail.UserID]
				if !ok {
					slog.ErrorContext(ctx, "user not found",
						"userID", detail.UserID,
						"orderID", order.ID.String(),
						"foodID", item.FoodID)
					continue
				}
				details = append(details, OrderItemDetail{User: userProfile, Quantity: detail.Quantity})
			}

			orderItems = append(orderItems, OrderItem{
				FoodName:     food.Name,
				FoodImageUrl: food.ImageURL,
				PriceAtOrder: item.PriceAtOrder,
				Details:      details,
			})
		}

		response = append(response, GetTodayOrdersResponse{
			PollID:     order.PollID,
			Buyer:      buyerProfile,
			TotalPrice: order.TotalPrice,
			OrderItems: orderItems,
		})
	}
	return response
}
