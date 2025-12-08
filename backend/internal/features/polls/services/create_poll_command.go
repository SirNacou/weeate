package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SirNacou/weeate/backend/internal/common/infrastructure/centrifugo"
	"github.com/SirNacou/weeate/backend/internal/features/auth"
	food_domain "github.com/SirNacou/weeate/backend/internal/features/foods/domain"
	"github.com/SirNacou/weeate/backend/internal/features/polls/domain"
	"github.com/gofrs/uuid/v5"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type CreatePollCommand struct {
	OrderDate        time.Time `json:"order_date" format:"date-time" doc:"Date for which the poll is being created in RFC3339 format"`
	ScheduledCloseAt time.Time `json:"scheduled_close_at" format:"date-time" doc:"Scheduled closing time for the poll in RFC3339 format"`
	FoodIDs          []string  `json:"food_ids" minItems:"1" doc:"List of food IDs to include in the poll"`
	Strategy         string    `json:"strategy" enum:"ORDER_MULTIPLE_ITEMS,ORDER_CONSENSUS_ITEM" doc:"Polling strategy to be used"`
}

type CreatePollResponse struct {
	PollID string `json:"poll_id" doc:"The ID of the created poll"`
}

type CreatePollCommandHandler struct {
	db         *gorm.DB
	centrifugo *centrifugo.CentrifugoClient
}

func NewCreatePollCommandHandler(db *gorm.DB, centrifugo *centrifugo.CentrifugoClient) *CreatePollCommandHandler {
	return &CreatePollCommandHandler{
		db:         db,
		centrifugo: centrifugo,
	}
}

func (h *CreatePollCommandHandler) Handle(ctx context.Context, req CreatePollCommand) (*CreatePollResponse, error) {
	user, ok := auth.UserFromContext(ctx)
	if !ok {
		return nil, auth.ErrUserNotFoundInContext
	}

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

	minScheduledTime := serverTime.Add(time.Minute * 15)
	if req.ScheduledCloseAt.UTC().Before(minScheduledTime.UTC()) {
		return nil, domain.ErrScheduledCloseAtTooSoon(minScheduledTime.Local())
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
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var existingPoll domain.Poll
		err := tx.Where("buyer_id = ? AND order_date = ?", user.ID, req.OrderDate).First(&existingPoll).Error
		if err == nil {
			if existingPoll.ClosedAt != nil {
				return domain.ErrClosedPollAlreadyExistsForOrderDate(req.OrderDate)
			}
			if err := tx.Unscoped().Delete(&existingPoll).Error; err != nil {
				return err
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		return tx.Create(poll).Error
	})
	if err != nil {
		return nil, err
	}

	_, err = h.centrifugo.PublishPublicPolls(ctx, &centrifugo.PublicPollsData{
		Type: centrifugo.PollCreated,
		Data: &struct {
			PollID uuid.UUID `json:"poll_id"`
		}{
			PollID: poll.ID,
		},
	})
	if err != nil {
		return nil, err
	}

	return &CreatePollResponse{
		PollID: poll.ID.String(),
	}, nil
}
