package services

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	auth_domain "github.com/SirNacou/weeate/backend/internal/features/auth/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type GetOrderedItemsQuery struct {
	Date time.Time `json:"date" format:"date-time" doc:"The date to get ordered items for"`
}

type GetOrderedItemsResponse struct {
	Items []OrderedItem `json:"items" nullable:"false"`
}

type OrderedItem struct {
	FoodID     string                  `json:"food_id"`
	FoodName   string                  `json:"food_name"`
	FoodURL    string                  `json:"food_url"`
	Quantity   int64                   `json:"quantity"`
	UnitPrice  int64                   `json:"unit_price"`
	TotalPrice int64                   `json:"total_price"`
	Buyer      auth_domain.UserProfile `json:"buyer"`
}

type GetOrderedItemsQueryHandler struct {
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewGetOrderedItemsQueryHandler(db *gorm.DB, supabaseService *auth.SupabaseService) *GetOrderedItemsQueryHandler {
	return &GetOrderedItemsQueryHandler{db: db, supabaseService: supabaseService}
}

func (h *GetOrderedItemsQueryHandler) Handle(ctx context.Context, query *GetOrderedItemsQuery) (*GetOrderedItemsResponse, error) {
	user, ok := auth_domain.UserFromContext(ctx)
	if !ok {
		return nil, auth_domain.ErrUserNotFoundInContext
	}

	type OrderItemRow struct {
		OrderID      uuid.UUID
		FoodID       uuid.UUID
		FoodName     string
		FoodURL      string
		Quantity     int64
		PriceAtOrder int64
		BuyerUserID  string
	}

	var rows []OrderItemRow
	err := h.db.Table("orders").
		Select(`
			orders.id as order_id,
			order_items.food_id,
			foods.name as food_name,
			foods.image_url as food_url,
			order_item_details.quantity,
			order_items.price_at_order,
			orders.buyer_user_id
		`).
		Joins("JOIN order_items ON order_items.order_id = orders.id").
		Joins("JOIN order_item_details ON order_item_details.order_item_id = order_items.id").
		Joins("LEFT JOIN foods ON foods.id = order_items.food_id").
		Where("orders.order_date = ?", query.Date).
		Where("order_item_details.user_id = ?", user.ID).
		Scan(&rows).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &GetOrderedItemsResponse{Items: []OrderedItem{}}, nil
		}
		return nil, err
	}

	buyerIDs := lo.UniqMap(rows, func(row OrderItemRow, _ int) string { return row.BuyerUserID })
	buyers, err := h.supabaseService.GetUserProfilesByIDs(buyerIDs...)
	if err != nil {
		return nil, err
	}
	buyersMap := lo.KeyBy(buyers, func(buyer auth_domain.UserProfile) string { return buyer.ID.String() })

	return &GetOrderedItemsResponse{
		Items: lo.Map(rows, func(row OrderItemRow, _ int) OrderedItem {
			return OrderedItem{
				FoodID:     row.FoodID.String(),
				FoodName:   row.FoodName,
				FoodURL:    row.FoodURL,
				Quantity:   row.Quantity,
				UnitPrice:  row.PriceAtOrder,
				TotalPrice: row.PriceAtOrder * row.Quantity,
				Buyer:      buyersMap[row.BuyerUserID],
			}
		}),
	}, nil
}
