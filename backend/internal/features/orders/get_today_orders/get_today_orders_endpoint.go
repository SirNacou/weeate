package get_today_orders

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/gofrs/uuid/v5"
)

type GetTodayOrdersRequest struct{}

type GetTodayOrdersResponse struct {
	PollID     uuid.UUID
	Buyer      auth.UserProfile
	TotalPrice int64
	OrderItems []OrderItem
}

type OrderItem struct {
	FoodName     string
	FoodImageUrl string
	PriceAtOrder int64
	Details      []OrderItemDetail
}

type OrderItemDetail struct {
	User     auth.UserProfile
	Quantity int
}

type GetTodayOrdersEndpoint struct {
	handler *GetTodayOrdersQueryHandler
}

func NewGetTodayOrdersEndpoint(handler *GetTodayOrdersQueryHandler) GetTodayOrdersEndpoint {
	return GetTodayOrdersEndpoint{
		handler: handler,
	}
}

func (e *GetTodayOrdersEndpoint) Handle(ctx context.Context, request *struct{ Body GetTodayOrdersRequest }) (*api.Response[[]GetTodayOrdersResponse], error) {
	orders, err := e.handler.Handle(ctx, &GetTodayOrdersQuery{})
	if err != nil {
		return nil, err
	}
	return api.NewResponse(orders), nil
}
