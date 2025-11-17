package orders

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
)

type OrdersEndpoint struct {
	getTodayOrdersQueryHandler *services.GetTodayOrdersQueryHandler
}

func NewOrdersEndpoint(g *services.GetTodayOrdersQueryHandler) *OrdersEndpoint {
	return &OrdersEndpoint{
		getTodayOrdersQueryHandler: g,
	}
}

func (e *OrdersEndpoint) getTodayOrders(ctx context.Context, request *struct{ Body services.GetTodayOrdersQuery }) (
	*api.Response[[]services.GetTodayOrdersResponse], error,
) {
	orders, err := e.getTodayOrdersQueryHandler.Handle(ctx, &request.Body)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(orders), nil
}
