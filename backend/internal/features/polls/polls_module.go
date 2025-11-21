package polls

import (
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterPollsModule(api huma.API, b *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) {
	getTodayPollsHandler := services.NewGetTodayPollsQueryHandler(db, supabaseService)
	createPollHandler := services.NewCreatePollCommandHandler(db)
	closePollHandler := services.NewClosePollCommandHandler(db, b)
	castVoteHandler := services.NewCastVoteCommandHandler(db)

	pollsEndpoint := NewPollsEndpoint(
		getTodayPollsHandler,
		createPollHandler,
		closePollHandler,
		castVoteHandler,
	)

	group := huma.NewGroup(api, "/polls")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Polls",
		Description: "Endpoints for managing polls",
	})

	huma.Get(group, "/today", pollsEndpoint.getTodayPolls)
	huma.Post(group, "/", pollsEndpoint.createPoll)
	huma.Post(group, "/{id}/close", pollsEndpoint.closePoll)
	huma.Post(group, "/{id}/vote", pollsEndpoint.castVote)
}
