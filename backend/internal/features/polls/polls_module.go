package polls

import (
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterPollsModule(api huma.API, b *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) {
	getActivePollsHandler := services.NewGetActivePollsQueryHandler(db, supabaseService)
	createPollHandler := services.NewCreatePollCommandHandler(db)
	closePollHandler := services.NewClosePollCommandHandler(db, b.EventBus)
	castVoteHandler := services.NewCastVoteCommandHandler(db)

	pollsEndpoint := NewPollsEndpoint(
		getActivePollsHandler,
		createPollHandler,
		closePollHandler,
		castVoteHandler,
	)

	group := huma.NewGroup(api, "/polls")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Polls",
		Description: "Endpoints for managing polls",
	})

	huma.Get(group, "/active", pollsEndpoint.getActivePolls)
	huma.Post(group, "/", pollsEndpoint.createPoll)
	huma.Post(group, "/close", pollsEndpoint.closePoll)
	huma.Post(group, "/vote", pollsEndpoint.castVote)
}
