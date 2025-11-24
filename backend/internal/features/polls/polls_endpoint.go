package polls

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
	"github.com/gofrs/uuid/v5"
)

type PollsEndpoint struct {
	getTodayPollsQueryHandler *services.GetTodayPollsQueryHandler
	createPollCommandHandler  *services.CreatePollCommandHandler
	closePollCommandHandler   *services.ClosePollCommandHandler
	castVoteCommandHandler    *services.CastVoteCommandHandler
}

func NewPollsEndpoint(
	get *services.GetTodayPollsQueryHandler,
	create *services.CreatePollCommandHandler,
	close *services.ClosePollCommandHandler,
	cast *services.CastVoteCommandHandler,
) *PollsEndpoint {
	return &PollsEndpoint{
		getTodayPollsQueryHandler: get,
		createPollCommandHandler:  create,
		closePollCommandHandler:   close,
		castVoteCommandHandler:    cast,
	}
}

func (e *PollsEndpoint) getTodayPolls(ctx context.Context, req *struct{}) (*api.Response[[]services.GetTodayPollsQueryResponse], error) {
	res, err := e.getTodayPollsQueryHandler.Handle(ctx, services.GetTodayPollsQuery{})
	if err != nil {
		return nil, err
	}
	return api.NewResponse(&res), nil
}

func (e *PollsEndpoint) createPoll(ctx context.Context, req *struct {
	Body services.CreatePollCommand
},
) (*api.Response[services.CreatePollResponse], error) {
	res, err := e.createPollCommandHandler.Handle(ctx, req.Body)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(res), nil
}

func (e *PollsEndpoint) closePoll(ctx context.Context, req *struct {
	ID uuid.UUID `path:"id" format:"uuid" doc:"The ID of the poll to be closed"`
},
) (*struct{}, error) {
	err := e.closePollCommandHandler.Handle(ctx, services.ClosePollCommand{PollID: req.ID})

	return nil, err
}

func (e *PollsEndpoint) castVote(ctx context.Context, req *struct {
	ID   uuid.UUID `path:"id" format:"uuid" doc:"The ID of the poll to cast a vote on"`
	Body struct {
		PollOptionID uuid.UUID `json:"poll_option_id" format:"uuid" doc:"The ID of the poll option to vote for"`
	}
},
) (*struct{}, error) {
	err := e.castVoteCommandHandler.Handle(ctx, services.CastVoteCommand{
		PollID:       req.ID,
		PollOptionID: req.Body.PollOptionID,
	})
	return nil, err
}
