package close_poll

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type ClosePollRequest struct {
	PollID uuid.UUID `json:"poll_id" format:"uuid" doc:"The ID of the poll to be closed"`
}

type ClosePollEndpoint struct {
	handler *ClosePollCommand
}

func NewClosePollEndpoint(h *ClosePollCommand) *ClosePollEndpoint {
	return &ClosePollEndpoint{
		handler: h,
	}
}

func (e *ClosePollEndpoint) Handle(ctx context.Context,
	req *struct {
		Body ClosePollRequest
	},
) (*struct{}, error) {
	err := e.handler.Handle(ctx, req.Body)
	return nil, err
}
