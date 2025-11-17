package create_poll

import (
	"context"
	"fmt"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/service"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreatePollCommand struct {
	FoodIDs []string
}

type CreatePollCommandHandler struct {
	db *gorm.DB
}

func NewCreatePollCommandHandler(db *gorm.DB) *CreatePollCommandHandler {
	return &CreatePollCommandHandler{
		db: db,
	}
}

func (h *CreatePollCommandHandler) Handle(ctx context.Context, req CreatePollRequest) (*CreatePollResponse, error) {
	user := ctx.Value(service.ContextKeyUser).(auth.User)

	foods, err := gorm.G[food_domain.Food](h.db).
		Where("id IN ?", req.FoodIDs).
		Where("user_id = ?", user.ID).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	foodsMap := lo.SliceToMap(foods, func(f food_domain.Food) (string, food_domain.Food) {
		return f.ID.String(), f
	})

	serverTime := time.Now()
	if req.OrderDate.Before(serverTime.Truncate(24 * time.Hour)) {
		return nil, domain.ErrOrderDateInPast
	}

	minScheduledTime := serverTime.Add(time.Hour)
	if req.ScheduledCloseAt.UTC().Before(minScheduledTime.UTC()) {
		return nil, domain.ErrScheduledCloseAtTooSoon(minScheduledTime)
	}

	var pollOptions []domain.PollOption
	for _, foodID := range req.FoodIDs {
		food, exists := foodsMap[foodID]
		if !exists {
			return nil, fmt.Errorf("food ID %s does not exist or does not belong to the user", foodID)
		}

		pollOptions = append(pollOptions, *domain.NewPollOption(uuid.UUID{}, food.ID, food.Price))
	}

	poll, err := domain.NewPoll(req.OrderDate,
		req.ScheduledCloseAt,
		domain.PollStrategy(req.Strategy),
		pollOptions,
	)
	if err != nil {
		return nil, err
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		gorm.G[domain.Poll](tx).Create(ctx, poll)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &CreatePollResponse{
		PollID: poll.ID.String(),
	}, nil
}
