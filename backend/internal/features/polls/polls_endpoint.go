package polls

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/polls/services"
)

type PollsEndpoint struct {
	getActivePollsQueryHandler *services.GetActivePollsQueryHandler
	createPollCommandHandler   *services.CreatePollCommandHandler
	closePollCommandHandler    *services.ClosePollCommandHandler
	castVoteCommandHandler     *services.CastVoteCommandHandler
}

func NewPollsEndpoint(
	get *services.GetActivePollsQueryHandler,
	create *services.CreatePollCommandHandler,
	close *services.ClosePollCommandHandler,
	cast *services.CastVoteCommandHandler,
) *PollsEndpoint {
	return &PollsEndpoint{
		getActivePollsQueryHandler: get,
		createPollCommandHandler:   create,
		closePollCommandHandler:    close,
		castVoteCommandHandler:     cast,
	}
}

func (e *PollsEndpoint) getActivePolls(ctx context.Context, req *struct{}) (*api.Response[[]services.GetActivePollsQueryResponse], error) {
	res, err := e.getActivePollsQueryHandler.Handle(ctx, services.GetActivePollsQuery{})
	if err != nil {
		return nil, err
	}
	return api.NewResponse(res), nil
}

func (e *PollsEndpoint) createPoll(ctx context.Context, req *struct {
	Body services.CreatePollCommand
},
) (*api.Response[services.CreatePollResponse], error) {
	res, err := e.createPollCommandHandler.Handle(ctx, req.Body)
	if err != nil {
		return nil, err
	}
	return api.NewResponse(*res), nil
}

func (e *PollsEndpoint) closePoll(ctx context.Context, req *struct{}) (*struct{}, error) {
	return nil, nil
}

func (e *PollsEndpoint) castVote(ctx context.Context, req *struct{}) (*struct{}, error) {
	return nil, nil
}
