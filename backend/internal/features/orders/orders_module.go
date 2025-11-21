package orders

import (
	"github.com/SirNacou/weeate/backend/internal/common/api/http"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterOrdersModule(api huma.API, bus *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) {
	group := huma.NewGroup(api, "/orders")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Orders",
		Description: "Endpoints for managing orders",
	})

	getTodayOrderHandler := services.NewGetTodayOrdersQueryHandler(db, supabaseService)
	ordersEndpoint := NewOrdersEndpoint(getTodayOrderHandler)
	huma.Get(group, "/today", http.Handle(ordersEndpoint.getTodayOrders))

	ordersEventHandler := NewOrdersEventHandler(
		services.NewCreateOrderCommandHandler(db, supabaseService),
	)

	bus.EventProcessor.AddHandlers(
		cqrs.NewEventHandler("orders.onPollClosed", ordersEventHandler.onPollClosed),
	)
}
