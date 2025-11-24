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
	endpoint           *OrdersEndpoint
	ordersEventHandler *OrdersEventHandler
	bus                *bus.Bus
}

func NewOrdersModule(bus *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) *OrdersModule {
	getTodayOrderHandler := services.NewGetTodayOrdersQueryHandler(db, supabaseService)
	ordersEndpoint := NewOrdersEndpoint(getTodayOrderHandler)

	ordersEventHandler := NewOrdersEventHandler(
		services.NewCreateOrderCommandHandler(db, supabaseService),
	)

	return &OrdersModule{
		endpoint:           ordersEndpoint,
		ordersEventHandler: ordersEventHandler,
		bus:                bus,
	}
}

func (m *OrdersModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/orders")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Orders",
		Description: "Endpoints for managing orders",
	})

	huma.Get(group, "/today", http.Handle(m.endpoint.getTodayOrders))

	m.bus.EventProcessor.AddHandlers(
		cqrs.NewEventHandler("orders.onPollClosed", m.ordersEventHandler.onPollClosed),
	)
}

func (m *OrdersModule) Jobs() []func(context.Context) {
	return nil
}
