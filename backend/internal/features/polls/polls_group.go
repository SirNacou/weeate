package polls

import (
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/get_active_polls"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterPollsGroup(db *gorm.DB, supabaseService *auth.SupabaseService, api huma.API) {
	getActivePollsHandler := get_active_polls.NewGetActivePollsQueryHandler(db, supabaseService)

	group := huma.NewGroup(api, "/polls")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Polls",
		Description: "Endpoints for managing polls",
	})

	getActivePollsEndpoint := get_active_polls.NewGetActivePollsEndpoint(getActivePollsHandler)
	huma.Get(group, "/active", getActivePollsEndpoint.Handle)
}
