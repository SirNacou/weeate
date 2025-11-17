package polls

import (
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/cast_vote"
	"github.com/SirNacou/weeate/backend/internal/features/polls/close_poll"
	"github.com/SirNacou/weeate/backend/internal/features/polls/create_poll"
	"github.com/SirNacou/weeate/backend/internal/features/polls/get_active_polls"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func RegisterPollsGroup(api huma.API, b *bus.Bus, db *gorm.DB, supabaseService *auth.SupabaseService) {
	getActivePollsHandler := get_active_polls.NewGetActivePollsQueryHandler(db, supabaseService)
	createPollHandler := create_poll.NewCreatePollCommandHandler(db)
	closePollHandler := close_poll.NewClosePollCommand(db)
	castVoteHandler := cast_vote.NewCastVoteCommandHandler(db)

	group := huma.NewGroup(api, "/polls")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Polls",
		Description: "Endpoints for managing polls",
	})

	getActivePollsEndpoint := get_active_polls.NewGetActivePollsEndpoint(getActivePollsHandler)
	huma.Get(group, "/active", getActivePollsEndpoint.Handle)

	createPollEndpoint := create_poll.NewCreatePollEndpoint(createPollHandler)
	huma.Post(group, "/", createPollEndpoint.Handle)

	closePollEndpoint := close_poll.NewClosePollEndpoint(closePollHandler)
	huma.Post(group, "/close", closePollEndpoint.Handle)

	castVoteEndpoint := cast_vote.NewCastVoteEndpoint(castVoteHandler)
	huma.Post(group, "/vote", castVoteEndpoint.Handle)
}
