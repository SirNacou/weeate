package orders

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/common/api/http"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
)

type OrdersEndpoint struct {
	getTodayOrdersQueryHandler   *services.GetTodayOrdersQueryHandler
	getOrderedItemsQueryHandler  *services.GetOrderedItemsQueryHandler
	getShoppingOrderQueryHandler *services.GetShoppingOrderQueryHandler
}

func NewOrdersEndpoint(g *services.GetTodayOrdersQueryHandler,
	o *services.GetOrderedItemsQueryHandler,
	s *services.GetShoppingOrderQueryHandler,
) *OrdersEndpoint {
	return &OrdersEndpoint{
		getTodayOrdersQueryHandler:   g,
		getOrderedItemsQueryHandler:  o,
		getShoppingOrderQueryHandler: s,
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

func (e *OrdersEndpoint) getOrderedItems(ctx context.Context, request *struct {
	Date time.Time `query:"date" format:"date-time" doc:"The date to get ordered items for"`
}) (
	*api.Response[services.GetOrderedItemsResponse], error,
) {
	// Round to start of day (midnight)
	date := time.Date(request.Date.Year(), request.Date.Month(), request.Date.Day(), 0, 0, 0, 0, request.Date.Location())

	orders, err := e.getOrderedItemsQueryHandler.
		Handle(ctx, &services.GetOrderedItemsQuery{
			Date: date,
		})
	if err != nil {
		return nil, http.MapError(err)
	}
	return api.NewResponse(orders), nil
}

func (e *OrdersEndpoint) getShoppingOrder(ctx context.Context, request *struct {
	Date time.Time `query:"date" format:"date-time" doc:"The date to get the shopping order for"`
}) (
	*api.Response[services.GetShoppingOrderResponse], error,
) {
	// Round to start of day (midnight)
	date := time.Date(request.Date.Year(), request.Date.Month(), request.Date.Day(), 0, 0, 0, 0, time.UTC)

	order, err := e.getShoppingOrderQueryHandler.
		Handle(ctx, &services.GetShoppingOrderQuery{Date: date})
	if err != nil {
		return nil, http.MapError(err)
	}
	return api.NewResponse(order), nil
}
