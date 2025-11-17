package cast_vote

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type CastVoteRequest struct {
	PollID       uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to vote in"`
	PollOptionID uuid.UUID `json:"poll_option_id" format:"uuid" doc:"The ID of the poll option to vote for"`
}

type CastVoteEndpoint struct {
	Handler *CastVoteCommandHandler
}

func NewCastVoteEndpoint(handler *CastVoteCommandHandler) *CastVoteEndpoint {
	return &CastVoteEndpoint{
		Handler: handler,
	}
}

func (e *CastVoteEndpoint) Handle(ctx context.Context,
	req *struct{ Body CastVoteRequest },
) (*struct{}, error) {
	err := e.Handler.Handle(ctx, req.Body)
	return nil, err
}
