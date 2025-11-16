package orders

import (
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/get_today_orders"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterOrdersGroup(api huma.API, db *gorm.DB, supabaseService *auth.SupabaseService) {
	group := huma.NewGroup(api, "/orders")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Orders",
		Description: "Endpoints for managing orders",
	})

	getTodayOrderHandler := get_today_orders.NewGetTodayOrdersQueryHandler(db, supabaseService)
	getTodayOrdersEndpoint := get_today_orders.NewGetTodayOrdersEndpoint(getTodayOrderHandler)
	huma.Get(group, "/today", getTodayOrdersEndpoint.Handle)
}
