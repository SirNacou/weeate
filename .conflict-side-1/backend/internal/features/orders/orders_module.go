package orders

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api/http"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type OrdersModule struct {
	endpoint     *OrdersEndpoint
	eventHandler *OrdersEventHandler
	bus          *bus.Bus
}

func NewOrdersModule(bus *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) *OrdersModule {
	getTodayOrderHandler := services.NewGetTodayOrdersQueryHandler(db, supabaseService)
	getOrderedItemsHandler := services.NewGetOrderedItemsQueryHandler(db, supabaseService)
	getShoppingOrderHandler := services.NewGetShoppingOrderQueryHandler(db, supabaseService)
	ordersEndpoint := NewOrdersEndpoint(getTodayOrderHandler, getOrderedItemsHandler, getShoppingOrderHandler)

	ordersEventHandler := NewOrdersEventHandler(
		services.NewCreateOrderCommandHandler(db, supabaseService),
	)

	return &OrdersModule{
		endpoint:     ordersEndpoint,
		eventHandler: ordersEventHandler,
		bus:          bus,
	}
}

func (m *OrdersModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/orders")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Orders",
		Description: "Endpoints for managing orders",
	})

	huma.Get(group, "/today", http.Handle(m.endpoint.getTodayOrders))
	huma.Get(api, "/ordered-items", http.Handle(m.endpoint.getOrderedItems))
	huma.Get(api, "/shopping-order", http.Handle(m.endpoint.getShoppingOrder))

	m.bus.EventProcessor.AddHandlers(
		cqrs.NewEventHandler("orders.onPollClosed", m.eventHandler.onPollClosed),
	)
}

func (m *OrdersModule) Jobs() []func(context.Context) {
	return nil
}
