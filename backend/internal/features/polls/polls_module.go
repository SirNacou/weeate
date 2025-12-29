package polls

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/bus"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/supabase"
	"github.com/SirNacou/weeate/backend/internal/features/polls/infrastructure/jobs"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jonboulle/clockwork"
	"gorm.io/gorm"
)

type PollsModule struct {
	scheduler *jobs.ClosePollScheduler
	endpoint  *PollsEndpoint
}

func NewPollsModule(b *bus.Bus, db *gorm.DB, supabaseService *supabase.SupabaseService, centrifugo *centrifugo.CentrifugoClient) *PollsModule {
	getTodayPollsHandler := services.NewGetTodayPollsQueryHandler(db, supabaseService)

	closePollHandler := services.NewClosePollCommandHandler(db, b, centrifugo)
	closePollScheduler := jobs.NewClosePollScheduler(closePollHandler, db, clockwork.NewRealClock())
	createPollHandler := services.NewCreatePollCommandHandler(db, centrifugo, closePollScheduler)
	castVoteHandler := services.NewCastVoteCommandHandler(db, centrifugo)

	pollsEndpoint := NewPollsEndpoint(
		getTodayPollsHandler,
		createPollHandler,
		closePollHandler,
		castVoteHandler,
	)

	return &PollsModule{
		scheduler: closePollScheduler,
		endpoint:  pollsEndpoint,
	}
}

func (m *PollsModule) RegisterAPI(api huma.API) {
	group := huma.NewGroup(api, "/polls")
	group.OpenAPI().Tags = append(group.OpenAPI().Tags, &huma.Tag{
		Name:        "Polls",
		Description: "Endpoints for managing polls",
	})

	huma.Get(group, "/today", m.endpoint.getTodayPolls)
	huma.Post(group, "/", m.endpoint.createPoll)
	huma.Post(group, "/{id}/close", m.endpoint.closePoll)
	huma.Post(group, "/{id}/vote", m.endpoint.castVote)
}

func (m *PollsModule) Jobs() []func(context.Context) {
	return []func(context.Context){
		m.scheduler.Start,
	}
}
