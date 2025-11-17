package create_poll

import (
	"context"
	"time"
)

type CreatePollRequest struct {
	OrderDate        time.Time `json:"order_date" doc:"Date for which the poll is being created in RFC3339 format"`
	ScheduledCloseAt time.Time `json:"scheduled_close_at" doc:"Scheduled closing time for the poll in RFC3339 format"`
	FoodIDs          []string  `json:"food_ids" minItems:"1" doc:"List of food IDs to include in the poll"`
	Strategy         string    `json:"strategy" enum:"ORDER_MULTIPLE_ITEMS,ORDER_CONSENSUS_ITEM" doc:"Polling strategy to be used"`
}

type CreatePollResponse struct {
	PollID string `json:"poll_id" doc:"The ID of the created poll"`
}

type CreatePollEndpoint struct {
	Handler *CreatePollCommandHandler
}

func NewCreatePollEndpoint(handler *CreatePollCommandHandler) *CreatePollEndpoint {
	return &CreatePollEndpoint{
		Handler: handler,
	}
}

func (e *CreatePollEndpoint) Handle(ctx context.Context,
	req *struct{ Body CreatePollRequest },
) (*CreatePollResponse, error) {
	return e.Handler.Handle(ctx, req.Body)
}
