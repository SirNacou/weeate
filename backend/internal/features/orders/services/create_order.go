package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/SirNacou/weeate/backend/internal/features/auth"
	"github.com/SirNacou/weeate/backend/internal/features/orders/domain"
	polls_domain "github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
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

	orderItems := make([]domain.OrderItem, 0)
	for _, option := range req.Results {
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
