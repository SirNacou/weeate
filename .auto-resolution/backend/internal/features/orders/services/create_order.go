package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/domain"
	polls_domain "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type OptionResult struct {
	FoodID          uuid.UUID
	PriceAtCreation int64
	Votes           []VoteResult
}

type VoteResult struct {
	UserID   string
	Quantity int64
}

type CreateOrderCommand struct {
	PollID    uuid.UUID
	BuyerID   string
	OrderDate time.Time
	Strategy  polls_domain.PollStrategy
	ClosedAt  time.Time
	Results   []OptionResult
}

type CreateOrderCommandHandler struct {
	db              *gorm.DB
	supabaseService *auth.SupabaseService
}

func NewCreateOrderCommandHandler(db *gorm.DB, supabaseService *auth.SupabaseService) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{
		db:              db,
		supabaseService: supabaseService,
	}
}

func (h *CreateOrderCommandHandler) Handle(ctx context.Context, req *CreateOrderCommand) error {
	if req.PollID.IsNil() {
		slog.WarnContext(ctx, "cannot create order without poll ID",
			"buyerID", req.BuyerID,
			"orderDate", req.OrderDate)
		return nil
	}

	totalVote := lo.SumBy(req.Results, func(r OptionResult) int {
		return len(r.Votes)
	})

	if len(req.Results) == 0 || totalVote == 0 {
		slog.InfoContext(ctx, "poll closed with no votes, skipping order creation",
			"pollID", req.PollID,
			"buyerID", req.BuyerID,
			"orderDate", req.OrderDate)
		return nil // Don't return error - this is a valid case
	}

	r, err := gorm.G[domain.Order](h.db).Where("order_date = ?", req.OrderDate).
		Where("buyer_user_id = ?", req.BuyerID).
		Count(ctx, "")
	if err != nil {
		return err
	}

	if r > 0 {
		slog.WarnContext(ctx, "order already exists for buyer on given date",
			"buyerID", req.BuyerID,
			"orderDate", req.OrderDate)
		return domain.ErrOrderAlreadyExists
	}

	if req.Strategy == polls_domain.OrderConsensus {
		totalVotes := lo.SumBy(req.Results, func(option OptionResult) int64 {
			return int64(len(option.Votes))
		})

		if totalVotes == 0 {
			slog.WarnContext(ctx, "cannot create order with consensus strategy with zero votes",
				"buyerID", req.BuyerID,
				"orderDate", req.OrderDate)
			return domain.ErrCannotCreateOrderWithZeroVotes
		}

		var winningOption *OptionResult
		for i, option := range req.Results {
			if winningOption == nil || len(option.Votes) > len(winningOption.Votes) {
				winningOption = &req.Results[i]
			}
		}

		// Collect all unique users from all options
		allVotes := make(map[string]VoteResult)
		for _, option := range req.Results {
			for _, vote := range option.Votes {
				// Keep the highest quantity if user voted multiple times
				if existing, exists := allVotes[vote.UserID]; !exists || vote.Quantity > existing.Quantity {
					allVotes[vote.UserID] = vote
				}
			}
		}

		// Transfer all votes to the winning option
		winningOption.Votes = lo.Values(allVotes)
		req.Results = []OptionResult{*winningOption}
	}

	orderItems := make([]domain.OrderItem, 0)
	for _, option := range req.Results {
		if len(option.Votes) == 0 {
			continue
		}
		orderItemDetails := make([]domain.OrderItemDetail, 0)
		for _, vote := range option.Votes {
			orderItemDetails = append(orderItemDetails, *domain.NewOrderItemDetail(vote.UserID, vote.Quantity))
		}

		orderItems = append(orderItems, *domain.NewOrderItem(option.FoodID, option.PriceAtCreation, orderItemDetails))
	}

	order, err := domain.NewOrder(req.PollID, req.BuyerID, req.OrderDate, req.Strategy, orderItems)
	if err != nil {
		return err
	}

	err = gorm.G[domain.Order](h.db).Create(ctx, order)

	return err
}
