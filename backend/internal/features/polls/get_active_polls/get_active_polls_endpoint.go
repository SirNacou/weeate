package get_active_polls

import (
	"context"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/api"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
)

type GetActivePollsQueryResponse struct {
	ID                string              `json:"id"`
	Creator           auth.UserProfile    `json:"creator"`
	ScheduledClosesAt time.Time           `json:"scheduled_closes_at"`
	FinalTotalPrice   int64               `json:"final_total_price"`
	Strategy          domain.PollStrategy `json:"strategy"`
	ClosedAt          *time.Time          `json:"closed_at"`
	PollOptions       []PollOption        `json:"poll_options"`
}

type PollOption struct {
	ID              string `json:"id"`
	Food            Food   `json:"food"`
	PriceAtCreation int64  `json:"price_at_creation"`
	Votes           []Vote `json:"votes"`
}

type Food struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Vote struct {
	Voter auth.UserProfile `json:"voter"`
}

type GetActivePollsRequest struct{}

type GetActivePollsEndpoint struct {
	handler GetActivePollsQueryHandler
}

func NewGetActivePollsEndpoint(handler GetActivePollsQueryHandler) *GetActivePollsEndpoint {
	return &GetActivePollsEndpoint{
		handler: handler,
	}
}

func (e *GetActivePollsEndpoint) Handle(ctx context.Context, r *struct {
	Body GetActivePollsRequest
},
) (*api.Response[[]GetActivePollsQueryResponse], error) {
	results, err := e.handler.Handle(ctx, GetActivePollsQuery{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]GetActivePollsQueryResponse]{
		Body: results,
	}, nil
}
