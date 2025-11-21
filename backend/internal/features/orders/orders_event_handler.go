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
	createOrderHandler *services.CreateOrderCommandHandler,
) *OrdersEventHandler {
	return &OrdersEventHandler{
		createOrderHandler: createOrderHandler,
	}
}

func (h *OrdersEventHandler) onPollClosed(ctx context.Context, event *events.PollClosedEvent) error {
	// Register the CreateOrderOnPollClosedHandler for PollClosedEvent
	return h.createOrderHandler.Handle(ctx, &services.CreateOrderCommand{})
}
