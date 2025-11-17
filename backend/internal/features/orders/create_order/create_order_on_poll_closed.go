package create_order

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/events"
)

type CreateOrderOnPollClosedHandler struct {
	handler *CreateOrderCommandHandler
}

func NewCreateOrderOnPollClosedHandler(h *CreateOrderCommandHandler) *CreateOrderOnPollClosedHandler {
	return &CreateOrderOnPollClosedHandler{
		handler: h,
	}
}

func (h *CreateOrderOnPollClosedHandler) Handle(ctx context.Context, event *events.PollClosedEvent) error {
	return h.handler.Handle(ctx, &CreateOrderCommand{
		// Map event fields to command fields here
	})
}
