package orders

import (
	"context"
	"log/slog"

	"github.com/SirNacou/weeate/backend/internal/common/events"
	"github.com/SirNacou/weeate/backend/internal/features/orders/services"
	polls_domain "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/samber/lo"
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
	slog.InfoContext(ctx, "received PollClosedEvent, processing order creation",
		"pollID", event.PollID,
		"buyerID", event.BuyerID,
		"orderDate", event.OrderDate,
		"strategy", event.Strategy,
		"optionCount", len(event.Results))

	err := h.createOrderHandler.Handle(ctx, &services.CreateOrderCommand{
		PollID:    event.PollID,
		BuyerID:   event.BuyerID,
		OrderDate: event.OrderDate,
		Strategy:  polls_domain.PollStrategy(event.Strategy),
		ClosedAt:  event.ClosedAt,
		Results: lo.Map(event.Results, func(r events.OptionResult, _ int) services.OptionResult {
			return services.OptionResult{
				FoodID:          r.FoodID,
				PriceAtCreation: r.PriceAtCreation,
				Votes: lo.Map(r.Votes, func(v events.VoteResult, _ int) services.VoteResult {
					return services.VoteResult{
						UserID:   v.UserID,
						Quantity: int64(v.Quantity),
					}
				}),
			}
		}),
	})

	if err != nil {
		slog.ErrorContext(ctx, "failed to process PollClosedEvent",
			"pollID", event.PollID,
			"buyerID", event.BuyerID,
			"error", err)
	}

	return err
}
