package orders

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/api/http"
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

func (e *OrdersEndpoint) getTodayOrders(ctx context.Context, request *struct{}) (
	*api.Response[[]services.GetTodayOrdersResponse], error,
) {
	orders, err := e.getTodayOrdersQueryHandler.Handle(ctx, nil)
	if err != nil {
		return nil, http.MapError(err)
	}
	return api.NewResponse(&orders), nil
}
