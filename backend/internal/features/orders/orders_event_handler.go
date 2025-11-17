package orders

import (
	"context"

	"github.com/SirNacou/weeate/backend/internal/common/events"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
)

type OrdersEventHandler struct {
	createOrderHandler *services.CreateOrderCommandHandler
}

func NewOrdersEventHandler(
	createOrderOnPollClosedHandler *services.CreateOrderCommandHandler,
) *OrdersEventHandler {
	return &OrdersEventHandler{
		createOrderHandler: createOrderOnPollClosedHandler,
	}
}

func (h *OrdersEventHandler) createOrderOnPollClosed(ctx context.Context, event *events.PollClosedEvent) error {
	// Register the CreateOrderOnPollClosedHandler for PollClosedEvent
	return h.createOrderHandler.Handle(ctx, &services.CreateOrderCommand{})
}
